import os
import difflib
import io
import hashlib
import docker
from provision_base import DockerProvisionBase

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
            "/usr/local/etc/php/conf.d/app.ini"
        )
        self.container.restart()

    def getVolumes(self):
        volumes = DockerProvisionBase.getVolumes(self, "/app")
        volumes[os.path.realpath(self.appConfig.appPath)] = {
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
        """ Generate unique id based on configuration. """
        hashStr = self.image
        hashStr += str(self.appConfig.getBuildFlavor())
        extensions = self.appConfig.getRuntime().get("extensions", [])
        extensions.sort()
        extensionConfigs = self.config.get("extensions", None)        
        for extension in extensions:
            if type(extension) is not str: continue # TODO accept dict with additional configs?
            extensionConfig = extensionConfigs.get(extension, {})
            if not extensionConfig: continue
            if extensionConfig.get("core", False): continue
            hashStr += extension
        return hashlib.sha256(hashStr).hexdigest()

    def generateNginxConfig(self):
        webConfig = self.appConfig.getWeb()
        locations = webConfig.get("locations", {})
        appNginxConf = ""
        def addFastCgi(scriptName = ""):
            if not scriptName: scriptName = "$fastcgi_script_name"
            conf = ""
            conf += "\t\t\tfastcgi_split_path_info ^(.+?\.php)(/.*)$;\n"
            conf += "\t\t\tfastcgi_pass %s:9000;\n" % (self.container.name)
            conf += "\t\t\tfastcgi_param SCRIPT_FILENAME $document_root%s;\n" % scriptName
            conf += "\t\t\tinclude fastcgi_params;\n"
            return conf

        for path in locations:
            appNginxConf += "location %s {\n" % path
            
            # root
            root = locations[path].get("root", "") or ""
            appNginxConf += "\t\troot \"%s\";\n" % (
                ("/app/%s" % (root.strip("/"))).rstrip("/")
            )

            # headers
            headers = locations[path].get("headers", {})
            for headerName in headers:
                appNginxConf += "\t\tadd_header %s %s;\n" % (
                    headerName,
                    headers[headerName]
                )

            # passthru
            passthru = locations[path].get("passthru", False)
            if passthru:
                if passthru == True: passthru = "/index.php"
                appNginxConf += "\t\tlocation ~ /%s {\n" % passthru.strip("/")
                appNginxConf += "\t\t\tallow all;\n"
                appNginxConf += addFastCgi(passthru)
                appNginxConf += "\t\t}\n"
                #appNginxConf += "\t\tlocation / {\n"
                appNginxConf += "\t\ttry_files $uri /%s$is_args$args;\n" % passthru.strip("/")
                #appNginxConf += "\t\t}\n"

            # scripts
            scripts = locations[path].get("scripts", False)
            appNginxConf += "\t\tlocation ~ [^/]\.php(/|$) {\n"
            if scripts:
                appNginxConf += addFastCgi()
                if passthru:
                    appNginxConf += "\t\t\tfastcgi_index %s;\n" % (passthru.lstrip("/"))
            else:
                appNginxConf += "\t\t\tdeny all;\n"
            appNginxConf += "\t\t}\n"

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
                    appNginxConf += "\t\tlocation ~ %s {\n" % (ruleRegex)

                    # allow
                    if not rule.get("allow", True):
                        appNginxConf += "\t\t\tdeny all;\n"
                    else:
                        appNginxConf += "\t\t\tallow all;\n"

                    # passthru
                    passthru = rule.get("passthru", False)
                    if passthru:
                        appNginxConf += addFastCgi(passthru)

                    # expires
                    expires = rule.get("expires", False)
                    if expires:
                        appNginxConf += "\t\t\texpires %s;\n" % expires

                    # headers
                    headers = rule.get("headers", {})
                    for headerName in headers:
                        appNginxConf += "\t\t\tadd_header %s %s;\n" % (
                            headerName,
                            headers[headerName]
                        )

                    # scripts
                    scripts = rule.get("scripts", False)
                    appNginxConf += "\t\t\tlocation ~ [^/]\.php(/|$) {\n"
                    if scripts:
                        appNginxConf += addFastCgi()
                        if passthru:
                            appNginxConf += "\t\t\t\tfastcgi_index %s;\n" % (passthru.lstrip("/"))
                    else:
                        appNginxConf += "\t\t\t\tdeny all;\n"
                    appNginxConf += "\t\t\t}\n"

                    appNginxConf += "\t\t}\n"

            appNginxConf += "\t}\n"
        return appNginxConf