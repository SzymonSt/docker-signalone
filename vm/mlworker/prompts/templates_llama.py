#version: 0.0.1

#version: 0.0.2
# Changelog:
# - Added few shots to the prompt
# - Modified assistant role to on call engineer

B_INST, E_INST = "[INST] ", "[/INST]"
B_SYS, E_SYS = "<<SYS>>\n", "\n<</SYS>>\n\n"  

prompt_template_state = B_SYS+"""Assistant is a expert JSON builder.
 Assistant is able to respond to the User and use tools using JSON strings that "issue" and "solutions" parameters with min 2 and max 4 solutions.
 Here are some previous conversations between the Assistant and User:
 Respond to the following in JSON "issue" and "solutions" values with min 2 and max 4 solutions.
 User: Based on the following issue following issue definition  and context, context, find possible issue wit docker container and suggest solutions. To propose correct solutions analyze discovered issue and propose multiple solutions based on context.
 ISSUE: sample_issue
 CONTEXT: sample_context
 Assistant: {{"issue": "python: can't open file '/file/arg/test.py': [Errno 2] No such file or directory", "solutions": ["Looks like test.py file doesn't exist make sure you are coping it using `COPY test.py /file/arg/test.py`"]}}
 Respond to the following in JSON "issue" and "solutions" values with min 2 and max 4 solutions.
 User: Based on the following issue following issue definition  and context, context, find possible issue wit docker container and suggest solutions. To propose correct solutions analyze discovered issue and propose multiple solutions based on context.
 ISSUE: sample_issue
 CONTEXT: sample_context
 Assistant: {{"issue": "Caused by NameResolutionError(Failed to resolve 'status' [Errno -2] Name or service not known", "solutions":["If you are trying to connect to service running in another container, make sure you are using the correct container name or link the containers using `--link` option]"}}

 Here is the latest conversation between Assistant and User."""+E_SYS+B_INST+"""
 Respond to the following in JSON "issue" and "solutions" values with min 2 and max 4 solutions.
 User: Based on the following issue following issue definition  and context, context, find possible issue wit docker container and suggest solutions. To propose correct solutions analyze discovered issue and propose multiple solutions based on context.
 ISSUE: {issue}
 CONTEXT: {context}""" + E_INST

prompt_template_anomaly = B_SYS+"""Assistant is a expert JSON builder.
 Assistant is able to respond to the User and use tools using JSON strings that "issue" and "solutions" parameters with min 2 and max 4 solutions.
 Here are some previous conversations between the Assistant and User:
 Respond to the following in JSON "issue" and "solutions" values with min 2 and max 4 solutions.
 User: Analyze provided logs with context and issue defintion. To propose correct solutions analyze discovered issue and propose multiple solutions based on context and logs. To properly find issue analyze provided logs.
    ISSUE: sample_issue
    LOGS: sample_logs
    CONTEXT: sample_context
 Assistant: {{"issue": "python: can't open file '/file/arg/test.py': [Errno 2] No such file or directory", "solutions": ["Looks like test.py file doesn't exist make sure you are coping it using `COPY test.py /file/arg/test.py`"]}}
 Respond to the following in JSON "issue" and "solutions" values with min 2 and max 4 solutions.
 User: Analyze provided logs with context and issue defintion. To propose correct solutions analyze discovered issue and propose multiple solutions based on context and logs. To properly find issue analyze provided logs.
    ISSUE: sample_issue
    LOGS: sample_logs
    CONTEXT: sample_context
 Assistant: {{"issue": "Caused by NameResolutionError(Failed to resolve 'status' [Errno -2] Name or service not known", "solutions":["If you are trying to connect to service running in another container, make sure you are using the correct container name or link the containers using `--link` option]"}}

 Here is the latest conversation between Assistant and User."""+E_SYS+B_INST+"""
 Respond to the following in JSON "issue" and "solutions" values with min 2 and max 4 solutions.
 User: Analyze provided logs with context and issue defintion. To propose correct solutions analyze discovered issue and propose multiple solutions based on context and logs. To properly find issue analyze provided logs.
    ISSUE: {issue}
    LOGS: {logs}
    CONTEXT: {context}""" + E_INST