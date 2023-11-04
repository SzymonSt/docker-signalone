import sqlite3

class IssuesManager():
    def __init__(self, db_path=None):
        if db_path:
            self.db = sqlite3.connect(db_path)
        else:
            self.db = sqlite3.connect("/data/issues.db")

    def insert_issues(self, issues):
        for issue in issues:
            issue['solutions'] = str(issue['solutions'])
            print(issue)
            self.db.execute("""INSERT INTO issues (
                id,
                container_id,
                issue_type,
                issue_severity,
                is_resolved,
                timestamp,
                issue,
                solutions) VALUES (
                    ?,
                    ?,
                    ?,
                    ?,
                    ?,
                    ?,
                    ?,
                    ?)""", [issue['id'], issue['container_id'], issue['issue_type'], issue['issue_severity'], issue['is_resolved'], issue['timestamp'], issue['issue'], issue['solutions']])
            self.db.commit()

    def drop_issues(self, issues):
        for issue in issues:
            self.db.execute("""DELETE FROM issues WHERE id = ?""", issue['id'])
            self.db.commit()

    def get_issues_by_container_id(self, container_id=None, filters=None, offset=0, limit=100):
        q = """SELECT * FROM issues"""
        parameters = []
        if container_id:
            q += """ WHERE container_id = ?"""
            parameters.append(container_id)
            if filters:
                q += """ AND """
        if filters:
            if not container_id:
                q += """ WHERE """
            for idx, filter in enumerate(filters):
                if filter['key'] == 'timestamp':
                    q += """{} >= ?""".format(filter['key'])
                    parameters.append(filter['value'])
                elif filter['key'] == 'end_timestamp':
                    q += """{} <= ?""".format("timestamp")
                    parameters.append(filter['value'])
                else:
                    q += """{} = ?""".format(filter['key'])
                    parameters.append(filter['value'])
                if idx != len(filters) - 1:
                    q += """ AND """
        q += " ORDER BY timestamp DESC LIMIT ? OFFSET ?"
        parameters.append(limit)
        parameters.append(offset)
        cur = self.db.execute(q, parameters)
        return cur.fetchall()
    
    def get_containers(self):
        cur = self.db.execute("""SELECT DISTINCT container_id FROM issues""")
        return cur.fetchall()