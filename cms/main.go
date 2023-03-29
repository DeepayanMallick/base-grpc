package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/go-chi/chi"

	user "github.com/DeepayanMallick/base-grpc/gunk/v1/usermgm/user"

	"google.golang.org/grpc"
)

const (
	ADDRESS  = "localhost:5000"
	homePath = "/"
)

type Server struct {
	Templates *template.Template
	User      user.UserServiceClient
}

type User struct {
	ID        string         `db:"id"`
	FirstName string         `db:"first_name"`
	LastName  string         `db:"last_name"`
	Username  string         `db:"username"`
	Email     string         `db:"email"`
	Password  string         `db:"password"`
	Status    int32          `db:"status"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
	DeletedAt sql.NullTime   `db:"deleted_at,omitempty"`
	DeletedBy sql.NullString `db:"deleted_by,omitempty"`
}

func main() {
	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect : %v", err)
	}

	defer conn.Close()

	s := &Server{
		Templates: &template.Template{},
		User:      user.NewUserServiceClient(conn),
	}
	s.ParseTemplates()

	r := chi.NewRouter()
	r.Get("/", s.GetUserHandler)
	r.Get("/create", s.CreateUserHandler)
	http.ListenAndServe(":8000", r)
}

func (s *Server) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res, err := s.User.GetUserByEmail(ctx, &user.GetUserByEmailRequest{
		Email: "john@graphlogic.com",
	})

	if err != nil {
		log.Println("Unable to get create user data")
		return
	}

	data := User{
		ID:        res.User.ID,
		FirstName: res.User.FirstName,
		LastName:  res.User.LastName,
		Username:  res.User.Username,
		Email:     res.User.Email,
	}

	s.parseHomeTemplate(w, data)
}

func (s *Server) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, err := s.User.CreateUser(ctx, &user.CreateUserRequest{
		User: &user.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@graphlogic.com",
			Username:  "john",
			Password:  "123456",
			Status:    1,
		},
	})
	if err != nil {
		log.Println("Unable to get create user data")
		return
	}

	http.Redirect(w, r, homePath, http.StatusSeeOther)
}

func (s Server) parseHomeTemplate(w http.ResponseWriter, data any) {
	t := s.Templates.Lookup("home.html")
	if t == nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) ParseTemplates() error {
	templates := template.New("graphlogic-template").Funcs(template.FuncMap{
		"globalfunc": func(n string) string {
			return ""
		},
	}).Funcs(sprig.FuncMap())
	newFS := os.DirFS("assets/template")
	tmpl := template.Must(templates.ParseFS(newFS, "*.html"))
	if tmpl == nil {
		log.Fatalln("unable to parse templates")
	}

	s.Templates = tmpl
	return nil
}
