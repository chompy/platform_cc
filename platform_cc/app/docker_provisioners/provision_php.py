from __future__ import absolute_import
import os
import difflib
import io
import hashlib
import base36
import docker
import sys
from .provision_base import DockerProvisionBase

class DockerProvision(DockerProvisionBase):

    """ Provision a PHP container. """

    def provision(self):
        # install extensions
        if self.logger:
            self.logger.logEvent(
                "Install extensions.",
                self.logIndent
            )

        # apt-get update
        self.container.exec_run(
            ["apt-get", "update"]
        )

        defaultExts = self.config.get("default_extensions", [])
        appExts = self.appConfig.getRuntime().get("extensions", [])

        extensions = defaultExts
        for extension in appExts:
            if type(extension) is not str: continue # TODO accept dict with additional configs?
            if extension in extensions: continue
            extensions.append(extension)

        extensionConfigs = self.config.get("extensions", {})
        for extensionName in extensions:
            if self.logger:
                self.logger.logEvent(
                    "%s." % (extensionName),
                    self.logIndent + 1
                )
            extensionConfig = extensionConfigs.get(extensionName, {})
            if not extensionConfig:
                if self.logger:
                    self.logger.logEvent(
                        "Extension not available.",
                        self.logIndent + 2
                    )
                continue
            if extensionConfig.get("core", False):
                if self.logger:
                    self.logger.logEvent(
                        "No additional configuration nessacary.",
                        self.logIndent + 2
                    )
                continue
            depCmdKey = difflib.get_close_matches(
                self.image,
                extensionConfig.keys(),
                1
            )
            if not depCmdKey:
                if self.logger:
                    self.logger.logEvent(
                        "Extension not available.",
                        self.logIndent + 2
                    )
                continue
            self.container.exec_run(
                ["sh", "-c", extensionConfig[depCmdKey[0]]]
            )
        # parent method
        DockerProvisionBase.provision(self)

    def runtime(self):
        DockerProvisionBase.runtime(self)
        if self.logger:
            self.logger.logEvent(
                "Copy php config from vars.",
                self.logIndent
            )
        phpConfig = self.appConfig.getVariables().get("php", {})
        phpIniOutput = ""
        for key, value in phpConfig.items():
            phpIniOutput += "%s = %s\n" % (key, value)
        self.copyStringToFile(
            phpIniOutput,
            "/usr/local/etc/php/conf.d/app2.ini"
        )
        self.container.restart()

    def getVolumes(self):
        volumes = DockerProvisionBase.getVolumes(self, "/app")
        appPath = os.path.realpath(self.appConfig.appPath)
        # hack for docker toolbox for windows, use unix path
        if sys.platform in ["msys", "win32"]:
            appPath = appPath.split(":")
            appPath = "/%s/%s" % (
                appPath[0].lower(),
                ("/".join(appPath[1].split("\\"))).lstrip("/")
            )
        volumes[appPath] = {
            "bind" : "/app",
            "mode" : "rw"
        }
        """appVolumeKey = "%s_%s_%s_app" % (
            DockerProvisionBase.DOCKER_VOLUME_NAME_PREFIX,
            self.appConfig.projectHash[:6],
            self.appConfig.getName()
        )
        try:
            self.dockerClient.volumes.get(appVolumeKey)
        except docker.errors.NotFound:
            self.dockerClient.volumes.create(appVolumeKey)
        volumes[appVolumeKey] = {
            "bind" : "/app",
            "mode" : "rw"
        }"""
        return volumes

    def getUid(self):
        # generate uid to be unique to the project/app (rather then to container) since
        # 'build' hooks are committed
        hashStr = self.image
        hashStr += self.appConfig.projectHash
        hashStr += self.appConfig.getName()
        return base36.dumps(
            int(
                hashlib.sha256(hashStr.encode("utf-8")).hexdigest(),
                16
            )
        )

    def generateNginxConfig(self):
        webConfig = self.appConfig.getWeb()
        locations = webConfig.get("locations", {})
        appNginxConf = ""
        def addFastCgi(scriptName = ""):
            if not scriptName: scriptName = "$fastcgi_script_name"
            conf = ""
            conf += "\t\t\t\tfastcgi_split_path_info ^(.+?\.php)(/.*)$;\n"
            conf += "\t\t\t\tfastcgi_pass %s:9000;\n" % (
                str(self.container.attrs.get("NetworkSettings", {}).get("IPAddress", None)).strip()
            )
            conf += "\t\t\t\tfastcgi_param SCRIPT_FILENAME $document_root%s;\n" % scriptName
            conf += "\t\t\t\tinclude fastcgi_params;\n"
            return conf

        for path in locations:
            appNginxConf += "\t\tlocation %s {\n" % path
            
            # root
            root = locations[path].get("root", "") or ""
            appNginxConf += "\t\t\troot \"%s\";\n" % (
                ("/app/%s" % (root.strip("/"))).rstrip("/")
            )

            # headers
            headers = locations[path].get("headers", {})
            for headerName in headers:
                appNginxConf += "\t\t\tadd_header %s %s;\n" % (
                    headerName,
                    headers[headerName]
                )

            # passthru
            passthru = locations[path].get("passthru", False)
            if passthru and not locations[path].get("scripts", False):
                if passthru == True: passthru = "/index.php"
                appNginxConf += "\t\t\tlocation ~ /%s {\n" % passthru.strip("/")
                appNginxConf += "\t\t\t\tallow all;\n"
                appNginxConf += addFastCgi(passthru)
                appNginxConf += "\t\t\t}\n"
                #appNginxConf += "\t\tlocation / {\n"
                appNginxConf += "\t\t\ttry_files $uri /%s$is_args$args;\n" % passthru.strip("/")
                #appNginxConf += "\t\t}\n"

            # scripts
            scripts = locations[path].get("scripts", False)
            appNginxConf += "\t\t\tlocation ~ [^/]\.php(/|$) {\n"
            if scripts:
                appNginxConf += addFastCgi()
                if passthru:
                    appNginxConf += "\t\t\t\tfastcgi_index %s;\n" % (passthru.lstrip("/"))
            else:
                appNginxConf += "\t\t\t\tdeny all;\n"
            appNginxConf += "\t\t\t}\n"

            # allow
            allow = locations[path].get("allow", False)
            # TODO!
            # allow = false should deny access when requesting a file that does exist but
            # does not match a rule

            # rules
            rules = locations[path].get("rules", {})
            if rules:
                for ruleRegex in rules:
                    rule = rules[ruleRegex]
                    appNginxConf += "\t\t\tlocation ~ %s {\n" % (ruleRegex)

                    # allow
                    if not rule.get("allow", True):
                        appNginxConf += "\t\t\t\tdeny all;\n"
                    else:
                        appNginxConf += "\t\t\t\tallow all;\n"

                    # passthru
                    passthru = rule.get("passthru", False)
                    if passthru:
                        appNginxConf += addFastCgi(passthru)

                    # expires
                    expires = rule.get("expires", False)
                    if expires:
                        appNginxConf += "\t\t\t\texpires %s;\n" % expires

                    # headers
                    headers = rule.get("headers", {})
                    for headerName in headers:
                        appNginxConf += "\t\t\t\tadd_header %s %s;\n" % (
                            headerName,
                            headers[headerName]
                        )

                    # scripts
                    scripts = rule.get("scripts", False)
                    appNginxConf += "\t\t\t\tlocation ~ [^/]\.php(/|$) {\n"
                    if scripts:
                        appNginxConf += addFastCgi()
                        if passthru:
                            appNginxConf += "\t\t\t\t\tfastcgi_index %s;\n" % (passthru.lstrip("/"))
                    else:
                        appNginxConf += "\t\t\t\t\tdeny all;\n"
                    appNginxConf += "\t\t\t\t}\n"

                    appNginxConf += "\t\t\t}\n"

            appNginxConf += "\t\t}\n"
        return appNginxConf