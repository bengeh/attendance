package main

import(
    "fmt"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "text/template"
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

func main(){
    
    http.HandleFunc("/show", Show)
    http.HandleFunc("/new", New)
    http.HandleFunc("/insert", Insert)
    http.ListenAndServe(":8080", nil)
}

func New(w http.ResponseWriter, r *http.Request){
    tmpl.ExecuteTemplate(w, "New", nil)
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
    }
    defer db.Close()
    http.Redirect(w, r, "show", 301)
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