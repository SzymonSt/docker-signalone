import os
import argparse
from subtools.source import SourceSubtool

def main():
    parsedArgs = parseArguments()
    if parsedArgs.subtool == 'source':
        source = SourceSubtool()
        source.execute(parsedArgs)

def parseArguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("--subtool", 
                        help="tool to interact with datasets or sources", 
                        choices=['dataset', 'source'])
    parser.add_argument("--action", help="action to perform on the tool", choices=[
        'scrape', 'push', 
        'generate', 'create',
        'merge'
        ])
    
    ## source scrape
    parser.add_argument("--source", 
                        help="source to scrape", 
                        choices=['stackoverflow', 'github'], 
                        default='stackoverflow')

    ## source push
    parser.add_argument("--log-file", 
                        help="path to log file")
    parser.add_argument("--format",
                        help="format of the log file", 
                        choices=['json', 'csv', 'plain'])
    
    ## source generate
    parser.add_argument("--base-logs-file", 
                        help="path to log file to base generation on")
    
    ## dataset create
    parser.add_argument("--sources-cat",
                        help="category of sources to create dataset from",
                        choices=['real', 'synthetic'],
                        default='real')
    parser.add_argument("--labels-cat",
                        help="category of labels to create for dataset",
                        choices=['real', 'synthetic'],
                        default='real')
    
    ## dataset merge
    parser.add_argument("--version",
                        help="version of datasets to merge", choices=calculateAllowedVersions())
    
    args = parser.parse_args()
    if args.subtool == 'source':
        if args.action == 'push':
            try:
                f = open(args.log_file, 'r')
                f.close()
            except FileNotFoundError:
                print("Log file not found")
                exit(1)
        elif args.action == 'generate':
            try:
                f = open(args.base_logs_file, 'r')
                f.close()
            except FileNotFoundError:
                print("Base logs file not found")
                exit(1)
        
    return args

def calculateAllowedVersions():
    versions = []
    ds_path = '../datasets'
    ds_files = os.listdir(ds_path)
    raw_versions = [file.split('_')[-1].replace('.csv', '') for file in ds_files]
    versions = list(set(raw_versions))
    return versions

if __name__ == '__main__':
    main()