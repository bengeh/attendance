package main

import(
    "fmt"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "text/template"
    "golang.org/x/crypto/bcrypt"
    "time"
)


func dbConn() (db *sql.DB){
    dbDriver := "mysql"
    dbUser := "root"
    dbName := "test_blog"
    db, err := sql.Open(dbDriver, dbUser + "@/" + dbName + "?parseTime=True")
    if err != nil{
        panic(err.Error())
    }
    return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

type Cookie struct {
    Name       string
    Value      string
    Path       string
    Domain     string
    Expires    time.Time
    RawExpires string

// MaxAge=0 means no 'Max-Age' attribute specified.
// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
// MaxAge>0 means Max-Age attribute present and given in seconds
    MaxAge   int
    Secure   bool
    HttpOnly bool
    Raw      string
    Unparsed []string // Raw text of unparsed attribute-value pairs
}
    
type Attendee struct {
    Id int 
    Name string
    Additional_pax int
    Food_choice string
}

type Total struct{
    Attendee
    Count int
    Steak_count int
    Salmon_count int
}

type Credentials struct{
    Username string `json:"username", db:"username"`
    Password string `json:"password", db:"password"`
}

func main(){
    http.HandleFunc("/", Home)
    http.HandleFunc("/login", Login)
    http.HandleFunc("/signup", Signup)
    http.HandleFunc("/show", Show)
    http.HandleFunc("/new", New)
    http.HandleFunc("/bad", Bad)
    http.HandleFunc("/insert", Insert)
    http.HandleFunc("/thanks", Thanks)
    http.ListenAndServe(":8080", nil)
}


func Signup(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    fmt.Println("Inside signup")
    // Parse and decode the request body into a new `Credentials` instance
    fmt.Println("Inside signup 2")
    name := r.FormValue("username")
    password := r.FormValue("password")
    // Salt and hash the password using the bcrypt algorithm
    // The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
    fmt.Println("Inside signup 4")
    // Next, insert the username, along with the hashed password into the database
    _ , err = db.Query("insert into users values (?, ?)", name, string(hashedPassword))
    if err != nil {
        // If there is any issue with inserting into the database, return a 500 error
        w.WriteHeader(http.StatusInternalServerError)
        return
    }else{
        http.Redirect(w, r, "New", 301)
    }
}


func Login(w http.ResponseWriter, r *http.Request){
    db := dbConn()
    // Get the existing entry present in the database for the given username
    fmt.Println("Inside login")
    username := r.FormValue("name")
    fmt.Print(username)
    password := r.FormValue("password")
    result, err:= db.Query("select password from test_blog.users where username = ?", username)
    var password1 string
    fmt.Println("Inside login 2")
    for result.Next(){
        err = result.Scan(&password1)
        fmt.Println(err)
        if err != nil {
            // If there is an issue with the database, return a 500 error
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
    }
    fmt.Println("Inside login 4")
    // Compare the stored hashed password, with the hashed version of the password that was received
    if err = bcrypt.CompareHashAndPassword([]byte(password1), []byte(password)); err != nil {
        // If the two passwords don't match, return a 401 status
        w.WriteHeader(http.StatusUnauthorized)
        http.Redirect(w, r, "bad", 301)
    }
    fmt.Println("Inside login 5")
    // If we reach this point, that means the users password was correct, and that they are authorized
    // The default 200 status is sent
    expiration := time.Now().Add(365 * 24 * time.Hour)
    cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
    http.SetCookie(w, &cookie)
    if username == "aaa"{
        http.Redirect(w, r, "show", 301)
    }else{
        http.Redirect(w, r, "new", 301)
    }
}


func Home(w http.ResponseWriter, r *http.Request){
    tmpl.ExecuteTemplate(w, "Home", nil)
}

func Bad(w http.ResponseWriter, r *http.Request){
    tmpl.ExecuteTemplate(w, "Bad", nil)
}

func Thanks(w http.ResponseWriter, r *http.Request){
    tmpl.ExecuteTemplate(w, "Thanks", nil)
}

func New(w http.ResponseWriter, r *http.Request){
    fmt.Print("Inside New")
    cookie, _ := r.Cookie("username")
    fmt.Print(cookie)
    tmpl.ExecuteTemplate(w, "New", cookie)
}

func Insert(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Inside insert")
    db := dbConn()
    if r.Method == "POST" {
        name := r.FormValue("name")
        add_pax := r.FormValue("add_pax")
        food_choice := r.FormValue("food_choice")
        fmt.Print(name)
        fmt.Print(add_pax)
        fmt.Print(food_choice)
        insForm, err := db.Prepare("INSERT INTO test_blog.wed_attendance(name, additional_pax, food_choice) VALUES(?,?,?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(name, add_pax, food_choice)
        fmt.Println("INSERT: Name: " + name + " | Additional pax: " + add_pax + " | Food choice: " + food_choice)
        if name == "aaa"{
            http.Redirect(w, r, "show", 301)
        }else{
            http.Redirect(w, r, "thanks", 301)
        }
    }
    defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request){
    fmt.Print("Hello World!")
    db := dbConn()
    rows, err := db.Query("SELECT * FROM test_blog.wed_attendance;")
    count := 0
    steak_count := 0
    salmon_count := 0
    attend := Total{}
    res := []Total{}
    for rows.Next(){
        var (
            id int
            name string
            additional_pax int
            food_choice string
        )
        
        err = rows.Scan(&id, &name, &additional_pax, &food_choice)
        if err != nil{
            panic(err)
        }
        attend.Id = id
        attend.Name = name
        attend.Additional_pax = additional_pax
        attend.Food_choice = food_choice
        
        if attend.Additional_pax == 1{
            count += 1
            attend.Count = count
        }
        if attend.Food_choice == "Steak"{
            steak_count += 1
            attend.Steak_count = steak_count
        }else if attend.Food_choice == "Salmon"{
            salmon_count += 1
            attend.Salmon_count = salmon_count
        }
        res = append(res, attend)
    }
    fmt.Println(res)
    tmpl.ExecuteTemplate(w, "Show", res)
    defer db.Close()
}