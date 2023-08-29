# PAN-USOM-XML2EDL

A Go program to fetch the USOM (TR-CERT) OSINT Feed in XML format and convert it to PANOS EDL format.

## Brief:

This program's existence is due to the following reasons:

* PANOS can only read EDLs from plain TXT files. It does not support other formats like XML, RSS, STIX, or TAXII. It is sensitive to data format and only accepts one value per line.
* PANOS has no regex capability for EDLs to parse different types. It cannot segregate different types like IP, domain, and URL from a single EDL source. Each type needs to be defined separately.
* PANOS has capacity limits for EDLs. Higher end platforms have a 250k URL limit, and lower end platforms have a 100k URL limit. Other limits exist for IPs and domains depending on the platform.
* USOM does not provide separate feeds for IP, domain, and URL IOC types, unlike other CERTs. They only provide a single source file with a mixture of each IOC type.
* USOM provides that single source in two formats: TXT and XML. Even in XML format, there is no metadata included for IOCs to identify their type. So regular expression operations are a must to identify types.
* USOM contains different and multiple syntaxes for each IOC type. They have no standard definitions for IOCs. Some sort of normalization is needed before use.
* USOM feed contains a huge number of IOCs. The number exceeds the PANOS capacity limits even for the highest plarforms. As they never remove a record from the list, it contains ancient and questionable records.
* Minemeld could be used to address some of these downsides formerly, but it is archived and not developed by PAN anymore. Also, it produces another link in the security toolchain on its own, which needs to be learned and maintained.

To overcome these challenges, PAN-USOM-XML2EDL needs to be used as a middleware. This program:

* Fetches USOM IOC Feed in XML format.
* Marshalls the XML data to match metadata for each IOCs,
* Filters IOCs by date using metadata information.
* Limits IOCs by count to stay within platform limits.
* Identifies and extracts different types of IOCs from a single feed.
* Parses multiple syntaxes and normalizes IOCs for each type.
* Compacts and de-duplicates recurring IOCs for each type.
* Segregates different types of IOCs and presents a seperate feed for each.

## Usage:

Pre-compiled binaries can be downloaded directly from the latest release (Link here: [Latest Release](https://github.com/enginy88/PAN-USOM-XML2EDL/releases/latest)). These binaries can readily be used on the systems for which they were compiled. Neither re-compiling any source code nor installing GO is not needed. In case there is no pre-compiled binary presented for your system, you can refer to the [Compilation](#compilation) section.

This program only requires the `appsett.env` file to determine which settings it will run with. Even if it can run with its default settings without any parameters set in the `appsett.env` file, this file must be present and accessible.

By default, the program searches for the `appsett.env` file in the working directory. The working directory can be altered by passing `-dir [PATH]` argument. The working directory is also determines where the generated EDL file(s) will be placed. If you need to change the output directory without changing the working directory, or if you have a scenario where you need to read the `appsett.env` file from one directory but create EDL files in a different directory, you can use the `-out [PATH]` argument for that purpose. The `-out [PATH]` argument overrides the working directory setting for the EDL file(s) that will be created. These options can be explored by passing the `-usage` argument to the program.

There are multiple settings controlled in environment variable format in the `appsett.env` file. These settings are explained in the [Settings](#settings) section.

How to schedule this program and how to serve generated EDL file(s) are not within the scope of this program or this documentation. Nevertheless, some hints are shared in the [Hints](#hints) section.

## Settings:

All setting options are provided with the sample `appsett.env` file, along with short descriptions and default values for each. Note that all lines are commented-out in the sample. If any of the options are to be used simply clear the comment token (`#`) and set the preferred setting.

```shell
# Basic & Optional Settings:
PANUSOMXML2EDL_FEED_URL={Enter URL of USOM XML OSINT Feed, Default: https://www.usom.gov.tr/url-list.xml}
PANUSOMXML2EDL_DAYS_OLD={Enter Max Day Old to Filter Entries, Default: 0 (Unlimited)}
PANUSOMXML2EDL_LIMIT_COUNT={Enter Limit Count to Filter Entries, Default: 0 (Unlimited)}
PANUSOMXML2EDL_NO_SORT={Enter Either TRUE or FALSE to Skip Sorting XML Entries, Default: FALSE}
PANUSOMXML2EDL_SINGLE_OUTPUT={Enter Either TRUE or FALSE to Create Single Output Without Any Parsing, Default: FALSE}

# Advanced & Expert Options:
PANUSOMXML2EDL_LOG_SILENCE={Enter Either TRUE or FALSE to Enable Verbose Log Suppression, Default: FALSE}
PANUSOMXML2EDL_LOG_UNPARSABLE={Enter Either TRUE or FALSE to Log Each Unparsable Record, Default: FALSE}
PANUSOMXML2EDL_SKIP_VERIFY={Enter Either TRUE or FALSE to Skip TLS Certificate Verification, Default: FALSE}
PANUSOMXML2EDL_FILE_TEST={Enter Either TRUE or FALSE to Run Offline Test from url-list.xml File, Default: FALSE}
```

<details>

<summary>Long explanation of each settings option: (Expand to view)</summary>

### Explanation of Settings:

**PANUSOMXML2EDL_FEED_URL** 

TYPE: ```String``` DEFAULT VALUE: ```https://www.usom.gov.tr/url-list.xml``` 

This setting controls which link the program will use to fetch the IOC Feed from USOM. Normally, there is no need to change this from the default value. However, it is implemented for a possible future scenario where USOM changes the link and  an admin cannot be able to quickly change it from the source code and then re-compile.

**PANUSOMXML2EDL_DAYS_OLD**

TYPE: ```Integer``` DEFAULT VALUE: ```0 (Unlimited)``` 

This setting is to filter records only specified days old. For instance, to filter only records added within 1 year, a value of 365 can be used.

**PANUSOMXML2EDL_LIMIT_COUNT**

TYPE: ```Integer``` DEFAULT VALUE: ```0 (Unlimited)``` 

This setting is to limit the output to a specified count of records. For instance, in a scenario where you have a device with a limit of 100k EDL, a value of 100,000 can be used to ensure it does not exceed the platform limit. Note that this limit applies during the internal process of the fetched records. So, the count of the output file will have this value of lines in Single-Output mode, but each output file will likely have fewer lines than the limit in normal (Multi-Output) mode due to separation and de-duplication operations.

**PANUSOMXML2EDL_NO_SORT**

TYPE: ```Boolean``` DEFAULT VALUE: ```FALSE``` 

Normally, USOM provides the XML file which contains records in appearing order date sorted. So for now, changing this option will have no effect. However, it is implemented for a possible future scenario where this behavior changes and there is a need to select records without sorting by date.


**PANUSOMXML2EDL_SINGLE_OUTPUT**

TYPE: ```Boolean``` DEFAULT VALUE: ```FALSE``` 

When set, the behavior of the program changes dramatically as it will not parse, process, separate, and de-duplicate the records and will create only a single output file named ```edl.txt```. When not set, which is the normal execution mode, the program parses and processes the records, compacts them by applying de-duplication, then separates each type of record and places them under files named ```edl-ip.txt```, ```edl-domain.txt``` & ```edl-url.txt```.

**PANUSOMXML2EDL_LOG_SILENCE**

TYPE: ```Boolean``` DEFAULT VALUE: ```FALSE``` 

This program has 4 levels of log types: Always, error, warning and log. When this option is set, it will suppress the warning and info level logs. Note that the level of always cannot be silenced. Also, the level of error is shown all the time as it means there is an unrecoverable failure during the operation and the program will be terminated without completing the jobs.

**PANUSOMXML2EDL_LOG_UNPARSABLE**

TYPE: ```Boolean``` DEFAULT VALUE: ```FALSE``` 

When set, this option will print all unparsable records in the feed. It may be useful to troubleshoot which lines are problematic.

**PANUSOMXML2EDL_SKIP_VERIFY**

TYPE: ```Boolean``` DEFAULT VALUE: ```FALSE``` 

Normally, USOM serves the XML file via the HTTPS protocol with a trusted TLS certificate. The default behavior is not to continue when the program encounters an untrusted certificate as it may indicate possible MITM attack, which tries to alter your block list. So, it is not advised to change this from the default value. However, it is implemented for possible use with in combination with a custom feed URL which may be served with an untrusted TLS certificate within your knowledge so that the program can proceed with an insecure connection.

**PANUSOMXML2EDL_FILE_TEST**

TYPE: ```Boolean``` DEFAULT VALUE: ```FALSE``` 

Normally, this program fetches the XML feed from either its default URL or a set custom URL. As of writing this document, the filesize of the feed exceeds 35MB so it may be inconvenient to download every time while testing. In addition to that, USOM may temporarily block you if you try to many requests as they have some sort of rate limiting mechanism. When this option set, the program looks for the ```url-list.xml``` file under working the directory and uses it for input. It may be useful for testing parsing capability.

</details>

## Hints:

When using this program under Unix-like OSes like Linux or macOS, Cron can be used to schedule periodic execution of the program. Here is an example of Crontab file:

```shell
# /etc/crontab
# To run the program every hour:
0 * * * * /path_to_binary/PAN-USOM-XML2EDL -dir /path_to_appsett.env_file/ -out /path_for_edl_files/
```
Under Windows OSes, the Task Scheduler tool can be used for the same purpose.

For serving the generated EDL file(s), Apache, Nginx, or IIS can be used to handle incoming HTTP(S) requests. For short-term testing purposes, http.server module of Python 3 can be used. Here is an example of how to run this:

```python
# Serves content of the current working directory:
python3 -m http.server 8080
```

## Compilation:

If none of the pre-compiled binaries covers your environment, you can choose to compile from source by yourself. Here are the instructions for that:

```shell
git clone https://github.com/enginy88/PAN-USOM-XML2EDL.git
cd PAN-USOM-XML2EDL
go mod init
go mod tidy
make local # To compile for your own environment.
make # To compile for pre-selected environments. 
```

**NOTE:** To compile from the source code, GO must be installed in the environment. However, it is not necessary for to run compiled binaries! Please consult the GO website for installation instructions: [Installing Go](https://go.dev/doc/install)

## Why Golang?

Because it (cross) compiles into machine code! You can directly run ready-to-go binaries on Windows, Linux, and macOS. No installation, no libraries, no dependencies, no prerequisites... Unlike Bash/PowerShell/Python it is not interpreted at runtime, which drastically reduces runtime overhead compared to scripting languages. The decision to use a compiled language makes it run lightning fast with lower memory usage. Also, due to the statically typed nature of the Go language, it is more error-proof against possible bugs/typos.