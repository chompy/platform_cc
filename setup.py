"""A setuptools based setup module.

See:
https://packaging.python.org/en/latest/distributing.html
https://github.com/pypa/sampleproject
"""

# Always prefer setuptools over distutils
from setuptools import setup, find_namespace_packages
# To use a consistent encoding
from codecs import open
from os import path

here = path.abspath(path.dirname(__file__))

# Get the long description from the README file
with open(path.join(here, 'README.md'), encoding='utf-8') as f:
    long_description = f.read()

setup(
    name='platform_cc',

    # Versions should comply with PEP440.  For a discussion on single-sourcing
    # the version across setup.py and the project code, see
    # https://packaging.python.org/en/latest/single_source_version.html
    version='0.3.14',

    description="Tool for provisioning apps with Docker based on Platform.sh's .platform.app.yaml spec.",
    long_description=long_description,

    # The project's main homepage.
    url='https://gitlab.com/contextualcode/platform_cc',

    # Author details
    author='Nathan Ogden @ Contextual Code',
    author_email='nathan@contextualcode.com',

    # Choose your license
    license='',

    # See https://pypi.python.org/pypi?%3Aaction=list_classifiers
    classifiers=[
    ],

    # What does your project relate to?
    keywords='development php docker platform.sh',

    # You can just specify the packages manually here if your project is
    # simple. Or you can use find_packages().
    packages=['platform_cc'] + find_namespace_packages(include=['platform_cc.*']),

    # Alternatively, if you want to distribute just a my_module.py, uncomment
    # this:
    #   py_modules=["my_module"],

    # List run-time dependencies here.  These will be installed by pip when
    # your project is installed. For an analysis of "install_requires" vs pip's
    # requirements files see:
    # https://packaging.python.org/en/latest/requirements.html
    install_requires=[
        'pyaml',
        'docker>=2.5',
        'cleo<0.7',
        'base36',
        'yamlordereddictloader',
        'terminaltables',
        'future',
        'dockerpty',
        'boto3',
        'nginx-config-builder',
        'cryptography'
    ],

    # List additional groups of dependencies here (e.g. development
    # dependencies). You can install these using the following syntax,
    # for example:
    # $ pip install -e .[dev,test]
    extras_require={
    },

    # If there are data files included in your packages that need to be
    # installed, specify them here.  If using Python 2.6 or less, then these
    # have to be included in MANIFEST.in as well.
    package_data={
        'platform_cc.core': ['data/*']
    },

    # Although 'package_data' is the preferred approach, in some case you may
    # need to place data files outside of your packages. See:
    # http://docs.python.org/3.4/distutils/setupscript.html#installing-additional-files # noqa
    # In this case, 'data_file' will be installed into '<sys.prefix>/my_data'
    data_files=[],

    # To provide executable scripts, use entry points in preference to the
    # "scripts" keyword. Entry points provide cross-platform support and allow
    # pip to create the appropriate form of executable for the target platform.
    entry_points={
        'console_scripts': [
            'platform_cc=platform_cc.commands:main',
        ],
    },
)
