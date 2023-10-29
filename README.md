# docker-signalone

## VM Setup

### api
    
```bash
    cd ./vm
```
Rename `./api/_env.example` to `.env`
```bash
    pip install -r requirements.txt
    python -m api.main
```
To access api documentation, go to `http://localhost:8000/docs`