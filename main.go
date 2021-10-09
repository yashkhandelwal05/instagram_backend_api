package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Model for User
type User struct{
	Id primitive.ObjectID  `json:"id" bson:"_id"`
	Name string				`json:"name" bson:"name"`
	Email string			`json:"email" bson:"email"`
	Password string			`json:"password" bson:"password"`
}
//Model for Post
type Post struct{
	Id primitive.ObjectID  `json:"id" bson:"_id"`
	Caption string			`json:"caption" bson:"caption"`
	ImageURL string			`json:"imageurl" bson:"imageurl"`
	Posted primitive.Timestamp `json:"date" bson:"Date()"`
}

type UserController struct{
	session *mgo.Session
}

func NewUserController(s *mgo.Session) *UserController{
	return &UserController{s}
}

//Get User using ID
func (uc UserController) GetUser (w http.ResponseWriter,r *http.Request,p httprouter.Params){
	id :=p.ByName("id");
	if(!bson.IsObjectIdHex(id)){
		w.WriteHeader(http.StatusNotFound)
	}
	oid :=bson.ObjectIdHex(id)
	u:=User{}
	if  err:=uc.session.DB("mongo-golang").C("users").FindId(oid).One(&u); err!=nil{
		w.WriteHeader(404)
		return
	}
	uj,err:=json.Marshal(u)
	if err!=nil{
		fmt.Println(err)
	}
	w.Header().Set("Content-type","application/json")
	w.WriteHeader(http.StatusOK)
	
	fmt.Printf("%s\n" , uj)

}

//Create User
func (uc UserController) CreateUser (w http.ResponseWriter,r *http.Request, _ httprouter.Params){
	u:=User{}
	json.NewDecoder(r.Body).Decode(&u)
	u.Id=primitive.NewObjectID()
	uc.session.DB("mongo-golang").C("users").Insert(u)
	uj,err:=json.Marshal(u)
	if err!=nil{
		fmt.Println(err)
	}
	w.Header().Set("Content-type","application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Printf("%s\n" , uj)
}

//Get Posts using ID
func (uc UserController) GetPosts (w http.ResponseWriter,r *http.Request,p httprouter.Params){
	id :=p.ByName("id");
	if(!bson.IsObjectIdHex(id)){
		w.WriteHeader(http.StatusNotFound)
	}
	oid :=bson.ObjectIdHex(id)
	u:=Post{}
	if  err:=uc.session.DB("mongo-golang").C("posts").FindId(oid).One(&u); err!=nil{
		w.WriteHeader(404)
		return
	}
	uj,err:=json.Marshal(u)
	if err!=nil{
		fmt.Println(err)
	}
	w.Header().Set("Content-type","application/json")
	w.WriteHeader(http.StatusOK)
	
	fmt.Printf("%s\n" , uj)
}

//Create Posts
func (uc UserController) CreatePosts (w http.ResponseWriter,r *http.Request, _ httprouter.Params){
	u:=Post{}
	json.NewDecoder(r.Body).Decode(&u)
	u.Id=primitive.NewObjectID()
	uc.session.DB("mongo-golang").C("posts").Insert(u)
	uj,err:=json.Marshal(u)
	if err!=nil{
		fmt.Println(err)
	}
	w.Header().Set("Content-type","application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Printf("%s\n" , uj)
}


func main() {
	r := httprouter.New()
	uc := NewUserController(getSession())
	r.GET("/user/:id",uc.GetUser)
	r.POST("/user",uc.CreateUser)
	r.GET("/posts/:id",uc.GetPosts)
	r.POST("/posts",uc.CreatePosts)
	//r.POST("/posts/users/:id",uc.GetPostsOfUser)
	http.ListenAndServe("localhost:9000",r)
}
func getSession() *mgo.Session{
	s,err:=mgo.Dial("mongodb://localhost:27017")
	if err != nil{
		panic(err);
	}
	return s;
}