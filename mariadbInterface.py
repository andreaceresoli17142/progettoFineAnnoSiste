import re
import mariadb, sys

class mariaDbInterface():
    def __init__(self, database):
        self.database = database

    def query( self, query ):

        try:
            conn = mariadb.connect(
                user="root",
                password="root",
                host="192.0.2.1",
                port=3306,
                database=self.database
            )
        except mariadb.Error as e:
            # print(f"Error connecting to MariaDB Platform: {e}")
            sys.exit(1)

        # Get Cursor
        cur = conn.cursor()
        allData = cur.execute(query)
        # cur.commit()
        conn.close()
        print(allData)
        return None
        if allData is None:
            return allData

        ret = []
        for row in allData:
            ret.append(row)
        return ret


dbi = mariaDbInterface("instanTex_db")

def main():
    print(dbi.query("SELECT * FROM Users;"))

if __name__ == "__main__":
    main()