package core

type ScriptData struct {
    Script		Script 	`redis:"script" json:"script"`
}

type Script struct {
    Name 		string  `redis:"name" json:"name"`
    Description string   `redis:"description" json:"description"`
	Fee		    []string `redis:"fee" json:"fee"`
}

type Fee struct {
    Name 		string  `redis:"name" json:"name"`
	Value		float64  `redis:"value" json:"value"`
}