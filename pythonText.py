from mariadbInterface import mariaDbInterface

dbi = mariaDbInterface("instanTex_db")

def main():
    print(dbi.query("SELECT * FROM Users;"))

if __name__ == "__main__":
    main()