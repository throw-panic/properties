package properties

import (
    "io"
    "bufio"
    "bytes"
    "unicode"
    "strconv"
)

type Properties struct {
    pairs map[string]string
}

func Load(reader io.Reader) (p *Properties, err error) {
    
    //  创建一个Properties对象
    p = new(Properties)
    p.pairs = make(map[string]string)
    
    //  创建一个扫描器
    scanner := bufio.NewScanner(reader)
    for scanner.Scan() {
        //  逐行读取
        line := scanner.Bytes()
        
        //  遇到空行
        if 0 == len(line) {
            continue
        }
        
        //  找到第一个非空白字符
        pos := bytes.IndexFunc(line, func(r rune) bool {
            return !unicode.IsSpace(r)
        })
        
        //  遇到空白行
        if -1 == pos {
            continue
        }
        
        //  遇到注释行
        if ('#' == line[pos]) || ('!' == line[pos]) {
            continue
        }
        
        //  找到第一个等号的位置
        end := bytes.Index(line[pos+1:], []byte{'='})
        
        //  没有=，说明该配置项只有key
        key := ""
        value := ""
        if -1 == end {
            key = string(bytes.TrimRightFunc(line[pos:], func(r rune) bool {
                return unicode.IsSpace(r)
            }))
        } else {
            key = string(bytes.TrimRightFunc(line[pos:pos+1+end], func(r rune) bool {
                return unicode.IsSpace(r)
            }))
            
            value = string(bytes.TrimSpace(line[pos+1+end+1:]))
        }
        
        p.pairs[key] = value
    }
    
    if err = scanner.Err(); nil != err {
        return nil, err
    }
    
    return p, nil
}

func (p*Properties) Get(key string) (value string, exist bool) {
    value, exist = p.pairs[key]
    return
}

func (p*Properties) StringDefault(key string, def string) string {
    value, ok := p.pairs[key]
    if ok {
        return value
    }
    
    return def
}

func (p*Properties) IntDefault(key string, def int64) int64 {
    value, ok := p.pairs[key]
    if ok {
        v, err := strconv.ParseInt(value, 10, 64)
        if nil != err {
            return def
        }
        
        return v
    }
    
    return def
}

func (p*Properties) FloatDefault(key string, def float64) float64 {
    value, ok := p.pairs[key]
    if ok {
        v, err := strconv.ParseFloat(value, 64)
        if nil != err {
            return def
        }
        
        return v
    }
    
    return def
}

func (p*Properties) BoolDefault(key string, def bool) bool {
    value, ok := p.pairs[key]
    if ok {
        v, err := strconv.ParseBool(value)
        if nil != err {
            return def
        }
        
        return v
    }
    
    return def
}

func (p*Properties) String(key string) string {
    return p.StringDefault(key, "")
}

func (p*Properties) Int(key string) int64 {
    return p.IntDefault(key, 0)
}

func (p*Properties) Float(key string) float64 {
    return p.FloatDefault(key, 0.0)
}

func (p*Properties) Bool(key string) bool {
    return p.BoolDefault(key, false)
}
