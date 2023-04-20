package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ratheeshkv179/TaskManger/persistence"
	"github.com/ratheeshkv179/TaskManger/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"sync"
)

var taskId = 0
var taskList map[string]types.Task
var mutex sync.Mutex

func getTaskId() string {
	return primitive.NewObjectID().Hex()
}

/*
func viewTasks(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Page Here viewTasks")
	var list []types.Task
	for _, v := range taskList {
		list = append(list, v)
	}
	if len(list) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	} else {
		data, err := json.Marshal(list)
		if err != nil {
			fmt.Errorf("Error: %#v\n", err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func viewTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Page Here viewTask")
	vars := mux.Vars(r)
	fmt.Println("Page Here viewTask", vars["id"])
	data, ok := taskList[vars["id"]]
	res, err := json.Marshal(data)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Internal server error", err.Error())))
	} else {
		if !ok {
			w.Write([]byte(fmt.Sprintf("Task with given Id %s does not exist", vars["id"])))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(res)
		}
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Page Here createTask")
	var data []byte
	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Errorf("Error: %#v\n", err.Error())
		w.Write([]byte(fmt.Sprintf("Failed to add the task, %s", err.Error())))
	} else {
		fmt.Println("Data", string(data))
		var task types.Task
		err := json.Unmarshal(data, &task)
		if err != nil {
			fmt.Errorf("Error: %#v\n", err.Error())
			w.Write([]byte(fmt.Sprintf("Failed to add the task, %s", err.Error())))
		} else {
			id := getTaskId()
			task.Id = id
			taskList[id] = task
			task.Status = types.CREATED
			data, err := json.Marshal(task)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("Internal server error", err.Error())))
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.Write(data)
			}
		}
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Page Here updateTask")
	vars := mux.Vars(r)
	fmt.Println("Page Here viewTask", vars["id"])
	oldData, ok := taskList[vars["id"]]
	if !ok {
		w.Write([]byte(fmt.Sprintf("Task Id %s does not exist", vars["id"])))
	} else {
		newData, err := io.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Internal server error %s", err.Error())))
		} else {
			var task types.Task
			err = json.Unmarshal(newData, &task)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("Internal server error %s", err.Error())))
			} else {
				task.Id = oldData.Id
				taskList[task.Id] = task
			}
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Page Here deleteTask")
	vars := mux.Vars(r)
	delete(taskList, vars["id"])
}
*/

func viewTasks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Page Here viewTasks")
	var list []types.TaskBase
	docs, err := dbClient.Get(DB_NAME, COLLECTION_NAME)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("DB query error %s", err.Error())))
	}
	for _, doc := range docs {
		b, err := bson.Marshal(doc)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Internal server error", err.Error())))
			return
		}
		var t types.TaskBase
		err = bson.Unmarshal(b, &t)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Internal server error", err.Error())))
			return
		}
		list = append(list, t)
	}
	if len(list) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	} else {
		data, err := json.Marshal(list)
		if err != nil {
			fmt.Errorf("Error: %#v\n", err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func viewTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Page Here viewTask")
	vars := mux.Vars(r)
	fmt.Println("Page Here viewTask", vars["id"])
	objectId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Invalid ObjectId %s, %s", vars["id"], err.Error())))
	} else {
		doc, err := dbClient.GetOne(DB_NAME, COLLECTION_NAME, bson.D{{"_id", objectId}})
		if err != nil {
			w.Write([]byte(fmt.Sprintf("DB query error %s", err.Error())))
		} else {
			b, err := bson.Marshal(doc)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("Internal server error", err.Error())))
				return
			}
			var t types.TaskBase
			err = bson.Unmarshal(b, &t)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("Internal server error", err.Error())))
				return
			}
			//data, ok := taskList[vars["id"]]
			res, err := json.Marshal(t)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("Internal server error", err.Error())))
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.Write(res)
			}
		}
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Page Here createTask")
	var data []byte
	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Errorf("Error: %#v\n", err.Error())
		w.Write([]byte(fmt.Sprintf("Failed to add the task, %s", err.Error())))
	} else {
		fmt.Println("Data", string(data))
		var task types.Task
		err := json.Unmarshal(data, &task)
		if err != nil {
			fmt.Errorf("Error: %#v\n", err.Error())
			w.Write([]byte(fmt.Sprintf("Failed to add the task, %s", err.Error())))
		} else {
			//id := getTaskId()
			task.Status = types.CREATED
			//taskList[id] = task
			updateDoc, err := dbClient.InsertOne(DB_NAME, COLLECTION_NAME, task)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("DB error while inserting doc", err.Error())))
			} else {
				b, err := bson.Marshal(updateDoc)
				if err != nil {
					w.Write([]byte(fmt.Sprintf("Internal server error", err.Error())))
				} else {
					var t types.TaskBase
					err = bson.Unmarshal(b, &t)
					if err != nil {
						w.Write([]byte(fmt.Sprintf("Internal server error", err.Error())))
					} else {

						json_bytes, err := json.Marshal(t)
						if err != nil {
							w.Write([]byte(fmt.Sprintf("Internal server error", err.Error())))
						} else {
							w.Header().Set("Content-Type", "application/json")
							w.Write(json_bytes)
						}

					}
				}
			}
		}
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Page Here updateTask")
	vars := mux.Vars(r)
	fmt.Println("Page Here viewTask", vars["id"])
	objectId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Invalid ObjectId %s, %s", vars["id"], err.Error())))
	} else {
		data, err := dbClient.GetOne(DB_NAME, COLLECTION_NAME, bson.D{{"_id", objectId}})
		if err != nil {
			w.Write([]byte(fmt.Sprintf("DB query error %s", err.Error())))
		} else {
			b, err := bson.Marshal(data)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("Internal server error", err.Error())))
				return
			}
			var oldTask types.TaskBase
			err = bson.Unmarshal(b, &oldTask)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("Internal server error", err.Error())))
				return
			}
			if err != nil {
				w.Write([]byte(fmt.Sprintf("Internal server error %s", err.Error())))
			} else {
				var newTask map[string]interface{}
				newData, err := io.ReadAll(r.Body)
				defer r.Body.Close()
				err = json.Unmarshal(newData, &newTask)
				if err != nil {
					w.Write([]byte(fmt.Sprintf("Internal server error %s", err.Error())))
				} else {
					//taskList[task.Id] = task
					objectId, err := primitive.ObjectIDFromHex(oldTask.Id)
					if err != nil {
						w.Write([]byte(fmt.Sprintf("Invalid ObjectId %s, %s", vars["id"], err.Error())))
					} else {
						var change primitive.E
						var changes []primitive.E
						for k, v := range newTask {
							change.Key = k
							change.Value = v
							changes = append(changes, change)
						}
						filter := bson.D{{"_id", objectId}}
						update := bson.D{{"$set", bson.D(changes)}}
						err := dbClient.UpSert(DB_NAME, COLLECTION_NAME, filter, update)
						if err != nil {
							w.Write([]byte(fmt.Sprintf("DB update error %s", err.Error())))
						}
					}
				}
			}
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Page Here deleteTask")
	vars := mux.Vars(r)
	objectId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Invalid ObjectId %s, %s", vars["id"], err.Error())))
	} else {
		err := dbClient.DeleteOne(DB_NAME, COLLECTION_NAME, bson.D{{"_id", objectId}})
		if err != nil {
			w.Write([]byte(fmt.Sprintf("DB delete error %s", err.Error())))
		}
	}
}

var dbClient persistence.MongoClient

const (
	DB_NAME         = "TaskManager"
	COLLECTION_NAME = "Tasks"
)

func main() {

	dbClient = persistence.MongoClient{}
	dbClient.Init("", "", "localhost", "27017")
	dbClient.Connect()

	taskList = make(map[string]types.Task)
	r := mux.NewRouter()
	r.HandleFunc("/tasks", viewTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", viewTask).Methods("GET")
	r.HandleFunc("/tasks", createTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", updateTask).Methods("PATCH")
	r.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	http.ListenAndServe(":80", r)
	fmt.Println("API server is shut down...")
	dbClient.Disconnect()
}
