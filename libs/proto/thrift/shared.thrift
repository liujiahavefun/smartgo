/**
 * This Thrift file can be included by other Thrift files that want to share
 * these definitions.
 */

//namespace cpp test2.shared
//namespace java test2.shared
//namespace php test2.shared
//namespace perl test2.shared
namespace go test2.shared

struct SharedStruct {
    1: i32 key
    2: string value
}

service SharedService {
    SharedStruct getStruct(1: i32 key)
}