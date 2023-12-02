import os
import json
import csv
import shutil

from helpers import helpers

class SourceSubtool(Subtool):
    def execute(self, args):
        if args.action == 'scrape':
            self.__scrape(args)
        elif args.action == 'push':
            self.__push(args)
        elif args.action == 'generate':
            self.__generate(args)
        else:
            print("Action not supported")
            exit(1)


    def __scrape(self, args):
        if args.source == 'stackoverflow':
            self.__scrapeStackOverflow()
        elif args.source == 'github':
            self.__scrapeGithub()

    def __push(self, args):
        if args.format == 'json':
            self.__pushJsonLogs(args.log_files)
        elif args.format == 'csv':
            self.__pushCsvLogs(args.log_files)
        elif args.format == 'plain':
            self.__pushPlainLogs(args.log_files)
    
    def __generate(self, args):
        self.__generateSyntheticSourceLogs(args.base_logs_file)

    def __scrapeStackOverflow(self):
        print("Scraping StackOverflow")
        raise NotImplementedError()
    
    def __scrapeGithub(self):
        print("Scraping Github")
        raise NotImplementedError()
    
    def __generateSyntheticSourceLogs(self, baseLogsFile):
        print("Generating synthetic source logs")
        raise NotImplementedError()
    
    def __pushJsonLogs(self, logFile):
        print("Pushing JSON logs to sources")
        key="body"
        with open(logFile, 'r') as f:
            logs = json.load(f.read())
            f.close()
        logs = [log[key] for log in logs]
        self.__pushLogs(logs)
           
    def __pushCsvLogs(self, logFile):
        print("Pushing CSV logs to sources")
        with open(logFile, 'r') as f:
            csv_reader = csv.reader(f)
            logs = [row for row in csv_reader]
            f.close()
        self.__pushLogs(logs)
    
    def __pushPlainLogs(self, logFile):
        print("Pushing plain logs to sources")
        with open(logFile, 'r') as f:
            logs = f.readlines()
            f.close()
        self.__pushLogs(logs)
    
    def __pushLogs(self, logs):
        batch_size=5
        log_batch_string = ""
        log_batch = []
        for log in logs:
            if batch_size <= 0:
                log_batch.append(log_batch_string)
                log_batch_string = ""
                batch_size = 5
            batch_size -= 1
            log_batch_string += log + "\n"
        if log_batch_string != "":
            log_batch.append(log_batch_string)
        file_template = 'collected_logs_v'
        current_version = self.__getNewSourceVersion(file_template)
        new_version = helpers.updateMinorVersion(current_version)
        current_file = file_template.replace('v', str(current_version))
        new_file  = file_template.replace('v', str(new_version))
        sources_path = '../sources/'
        shutil.copyfile(sources_path + current_file, sources_path + new_file)
        with open(sources_path + new_file, 'a', newline='') as f:
            csv_writer = csv.writer(f)
            for log in log_batch:
                csv_writer.writerow(log)
    
    def __getNewSourceVersion(self, source_file_template):
        files = os.listdir('../sources')
        files_matching_template = [file for file in files if file.startswith(source_file_template)].sort()
        if len(files_matching_template) == 0:
            return "v0.0.1"
        else:
            return files_matching_template[-1].split('_')[-1].replace('.csv', '')