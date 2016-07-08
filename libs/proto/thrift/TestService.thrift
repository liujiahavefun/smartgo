namespace go test.demo
namespace php test.demo

/**
 * 结构体定义
 */
struct Article{
    1: i32 id, 
    2: string title,
    3: string content,
    4: string author,
}

const map<string,string> MAPCONSTANT = {'hello':'world', 'goodnight':'moon'}

service test2Thrift {        
    list<string> CallBack(1:i64 callTime, 2:string name, 3:map<string, string> paramMap),
    void put(1: Article newArticle),
}