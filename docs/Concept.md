# DCX Base Concept

DCX consists of 2 parts, the orchestrator and the executor. 

POST /session API is used to submit code for execution with request body: 
Files stores the filesystem of the project, based from the root project directory.
Example:
```json
{
    "Files": {
        "main.py": {
            "encoding": "utf-8",
            "content": "a = input()\nprint(a)"
        },
        "requirements.txt": {
            "encoding": "utf-8",
            "content": "requests"
        }
    },
    "Language": "python",
    "Tests": [
        {
            "input": "Hello, World!",
            "output": "Hello, World!"
        }
    ]
}
```
Then it returns a session_id which is generated using UUID.

GET /status/:session_id returns the status of the session, with the following status:
- pending
- running
- done
- error

When pending, it means the session is waiting for the executor to pick it up.
When running, it means the session is being executed.
When done, it means the session is finished and the result is available.
When error, it means the session is finished and the result is error.

When done, the result is available in the response body:
```json
{
    "Status": "done",
    "Results": [
        {
            "Input": "1+1",
            "ActualOutput": "2",
            "ExpectedOutput": "2",
            "Status": "success"
        }
    ]
}
```

POST /submit
Request body:
```json
{
    "Files": {
        "main.py": {
            "Encoding": "utf-8",
            "Content": "a = input()\nprint(a)"
        },
        "requirements.txt": {
            "Encoding": "utf-8",
            "Content": "requests"
        }
    },
    "Image": "python:3.10",
    "Tests": [
        {
            "Input": "Hello, World!",
            "Output": "Hello, World!"
        }
    ],
    "Configurations": {
        "Strictness": "strict",
        "Timeout": 10,
        "MemoryLimit": 1024
    }
}
```
And in there, it stores to a local map and returns a ExecutionGroupID and each ExecutionID, and it's current status.
