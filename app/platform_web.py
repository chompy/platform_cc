import os
import yaml
from platform_config import PlatformConfig
from platform_docker import PlatformDocker
from app.platform_utils import log_stdout, print_stdout

class PlatformWeb:

    """ Provide web access to app via nginx docker container. """

    WEB_DOCKER_IMAGE = "nginx:1.13"

    def __init__(self, app):
        self.app = app
        self.docker = PlatformDocker(
            self.app.config,
            "web",
            self.WEB_DOCKER_IMAGE
        )

    def generateNginxConfig(self):
        """ Generate nginx config file for application. """
        webConfig = self.app.config.getWeb()
        locations = webConfig.get("locations", {})

        webProvisionConfig = self.docker.getProvisioner().config
        baseNginxConfig = webProvisionConfig.get("nginx.conf", "")
        appNginxConf = ""

        def addFastCgi(scriptName = ""):
            if not scriptName: scriptName = "$fastcgi_script_name"
            conf = ""
            conf += "\t\t\tfastcgi_split_path_info ^(.+?\.php)(/.*)$;\n"
            conf += "\t\t\tfastcgi_pass %s:9000;\n" % (self.app.docker.containerId)
            conf += "\t\t\tfastcgi_param SCRIPT_FILENAME $document_root%s;\n" % scriptName
            conf += "\t\t\tinclude fastcgi_params;\n"
            return conf

        for path in locations:
            appNginxConf += "location %s {\n" % path
            
            # root
            appNginxConf += "\t\troot \"%s\";\n" % (
                "/app/%s" % (locations[path].get("root", "").strip("/"))
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
                appNginxConf += "\t\tlocation ~ /%s {\n" % passthru.strip("/")
                appNginxConf += "\t\t\tallow all;\n"
                appNginxConf += addFastCgi(passthru)
                appNginxConf += "\t\t}\n"
                appNginxConf += "\t\tlocation / {\n"
                appNginxConf += "\t\t\ttry_files $uri /index.php$is_args$args;\n"
                appNginxConf += "\t\t}\n"

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

        return baseNginxConfig.replace("{{APP_WEB}}", appNginxConf)

    def start(self):
        """ Start app web handler. """
        print_stdout("> Starting web handler for '%s' app." % self.app.config.getName())
        self.docker.start()
        self.docker.getProvisioner().copyStringToFile(
            self.generateNginxConfig(),
            "/etc/nginx/nginx.conf"
        )
        self.docker.getContainer().restart()

    def stop(self):
        """ Stop app web handler. """
        print_stdout("> Stopping web handler for '%s' app." % self.app.config.getName())
        self.docker.stop()