# Go Configuration Reader

Allows to read configurations from different sources and combine all together in one structure.

## Basic example

```Go
type Config struct {
    Field_1 int                `env:"f1,required" def:"10"`
    Field_2 []string           `env:"f2,append" def:"a|b" sep:"|"`
    Field_3 [3]time.Time       `env:"f3" def:"now,now"`
    Field_4 map[string]float64 `env:"f4" def:"a?1.1,b?2.2" sep2:"?"`
    Field_5 SubConfig          `env:"f5"`
    Fiend_6 string             `env:"GOPATH"`
    Fiend_7 *float64           `env:"-"`
    Field_8 []*int8            `env:"f8" def:"*nil,1,2,3,4,5"`
    FIell_9 [2]*float32        `env:"f9,required" def:"1.23"`
}
type SubConfig struct {
    Field_1 *uint32            `env:"sf1" def:"*nil"`
}

func main() {
    config := &Config{}
    cr := goconf.NewConfigReader().
        RewriteValues(true).
        AddEnvironment().
        AddFile(".env").
        AddFile("config.ini").
        AddString("f5.sf1 = 10", goconf.FtEnv, "config 1").
        AddFile("config.json").
        EnsureHasNoErrors()
    err := cr.ReadConfig(config)
    if err != nil {
        panic(err)
    }

    ...
}
```

## Suppported field types
+ all numbers, string, bool, Time, Duration,  
+ pointers, slices, slices of pointers, arrays, arrays of pointers for all these types,
+ maps for these types but the key is always string,
+ substructs.
+ More types in future.

## Supported file types
+ env
+ ini
+ json  
+ ~~yaml~~ (later)

## Public methods

+ `NewConfigReader(files ...string)` - creates new instance of the configuration reader. You can add there all the config files paths. **Source order is important**!

+ `AddFile(file string)` - add the config file path as a configuration source.

+ `AddEnvironment()` - use environment variables as a configuration source.

+ `AddString(values string, formatType formatType, name string)` - add configuration source as a string.

+ `WithOptions(options ConfigOptions)` - add configuration reader parsing options.

+ `RewriteValues(rewrite bool)` - one of the options. Default is `true`. If different sources have different for the same key name defines weather will be used the first found (`false`) or the last one (`true`). Doesn't work for collections.

+ `WithParser(envName string, parser Parser)` - specify parser function for the specific structure.

+ `EnsureHasNoErrors()` - checks data before parsing and panics if wrong sources where added.

+ `GetErrors()` - checks data before parsing and returns errors if wrong sources where added.

+ `ReadConfig(userConfig interface{})` - reads configuration sources.

## Supported tags and options

+ `env` - configuration name for the field, case insensitive. If not set field name will be used. If `-`, field will be ignored. You can use any symbols but `.`.  
    **Env options** (comma separated)
    * `required` - if set, value in any source or default value must be specified. If no any value was found returns error. Collection must contain at least one item.
    * `append` - for collections only. If set, you'll get all the values from all sources.
    * `useparser` - if set, uses specified parser (works for any type).
+ `def` - default value. Can be used for any field type but structure without `useparser` option.
+ `sep` - separator for collections. Default is `,`.
+ `sep2` - separator between key and value in maps. Default is `:`.

Default built-in values, can be used in sources and `def` tag.  
+ `*nil` - sets nil if it's possible for the field (pointer, slice, map, item of collection of pointers).
+ `now` - sets current time for `Time` field.

If you have `1,2,3` for the array field of 5 `[5]int` you'll get `[1,2,3,0,0]`, but `1,2,3,4,5,6` will return an error.


## More examples

### Parser

```Go
type Config struct {
    Field_1 SubConfig  `env:"f1,useparser" def:"f1_1"`
    Field_2 SubConfig  `env:"f2,useparser" def:"f2_1"`
}
type SubConfig struct {
    Field_1 string
    Field_2 int
}

func main() {
    config := &Config{}
    cr := goconf.NewConfigReader().
        WithParser("f1", ParseSubConfig).
        WithParser("f2", ParseSubConfig).
        AddString("f1 = f1_2", goconf.FtEnv, "config 1").
        EnsureHasNoErrors()
    err := cr.ReadConfig(config)
    if err != nil {
        panic(err)
    }

    ...
}

func ParseSubConfig(s string) (interface{}, error) {
    split := strings.Split(s, "_")
    i, err := strconv.Atoi(split[1])
    if err != nil {
        return nil, err
    }

    return SubConfig{Field_1: split[0], Field_2: i}, nil
}
```

Here for the Field_1 will be parsed value `f1_2` from the string source, for the FIeld_2 will be parsed value `f2_1` from the default tag.

## Sources formats

### Env file
```ini
; comment

key_1=value
key_2 = "value"
key_3 = 'va"lu"e' ; returns va"lu"e

sub.key = value ; key in sub-structure

key_4 = 1,2,3 ; if field is a slice or an array, contains multiple values split by separator
key_4[] = 4 ; adds another one

key_5[] = value_1 ; only for arrays and slices. Contains single value
key_5[] = value_2 ; as many values as you set
key_5 = value_3,value_4 ; adds two more

key_6[k1] = v1 ; only for maps
key_6[k2] = v2

; key = value - this line will be ignored

key[_]7 = value ; just a regular key with the name "key[_]7"

key_8 = multy\
line\
; value. Returns "multy\nline\n", doesn't work with quoted lines
```

Keys cannot contain spaces.  
Values can contain spaces but quoted only.  
No one escape symbol is supported but quote which was used for quotion.

### Ini file
```ini
; comment
# another comment

# [ROOT] - default section
key 1 = value 1 # can contain spaces
key_2="value"
key_3 = 'va"lu"e' ; returns va"lu"e

[sub]
sub_key_1 = value ; values for the key in sub-structure
    [.sub 2] # if starts from `.`, subsection for previous section, can contain spaces
    sub_sub_key_1 = value
[sub.sub 2] # same with previous one
sub_sub_key_2 = value

; sub.key - key cannot contain `.`

[ rOOt ] # back to root section, these spaces will be ignored
key_4 = 1,2,3 ; if field is a slice or an array, contains multiple values split by separator
key_4[] = 4 ; adds another one

key_5[] = value_1 ; only for arrays and slices. Contains single value
key_5[] = value_2 ; as many values as you set
key_5 = value_3,value_4 ; adds two more

key_6[k1] = v1 ; only for maps
key_6[k2] = v2

; key = value - this line will be ignored

key[_]7 = value ; just a regular key with the name "key[_]7"

key_8 = multy\
line\
; value. Returns "multy\nline\n", doesn't work with quoted lines
```

Keys can contain anything but `.`.  
Sections, keys and values can contain spaces.  
No one escape symbol is supported but quote which was used for quotion.

### Json file
```javascript
// Comment
{
    "key_1":"value","key_2":"va\"lue",
    "/**/key_2[]": "value",  //possible key
    /*"key":"value"
        these lines will be ignored*/
    
    "key_3"//:
    :/*false*/
    true//,
    ,   // yes, that is possible, bool value

    "key_4": -123,  // number value
    "key_5": null, // null value
    "key_6": NuLl, // still null value, key words are case insensitive

    "key_7": [1,2,3,],   // array, trailing commas are possible anywhere
    
    "key_8": {  //sub-structure with two bool fields
        "k1": TRUE,
        "k2": FALSE,
    },

    "key_9": { //looks the same but it's for the map
        "k1": 1.1,
        "k2": 2.2
    },

    "key_10": { // structure with an array field
        "k1": ["a", "b"]
    },

    // "key_11": [1, "2"] - different types not allowed
}
/* some
comments here */
```

The root is always object!  
Keys can contain anything but `.`.  
No one escape symbol is supported but `\"`.

