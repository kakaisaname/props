package props

import (
    "unicode"
    "bufio"
    "bytes"
    "io"
    "os"
    "fmt"
    log "github.com/sirupsen/logrus"
)

const (
    TIME_S  = "S"
    TIME_MS = "MS"
)

type Properties struct {
    MapProperties
    //values map[string]string
}

func NewProperties() *Properties {
    p := &Properties{
        //values: make(map[string]string),
    }
    p.values = make(map[string]string)
    return p
}

// Read creates a new property set and fills it with the contents of a file.
// See Load for the supported file format.
func readProperties(r io.Reader) (*Properties, error) {
    p := NewProperties()
    err := p.Load(r)
    if err != nil {
        log.Error(err)
        return nil, err
    }
    return p, nil
}

func ReadPropertyFile(f string) (*Properties, error) {

    file, err := os.Open(f)
    defer file.Close()

    if err != nil {
        d, _ := os.Getwd()
        log.WithField("error", err.Error()).Info("read file: ", d, "  ", f)
        return nil, err
    }
    return readProperties(file)
}

func (p *Properties) Load(r io.Reader) error {

    buf := bufio.NewReader(r)
    for {
        line, _, err := buf.ReadLine()
        if err != nil {
            if err == io.EOF {
                return nil
            } else {
                return err
            }
        }

        line = bytes.TrimSpace(line)
        lenLine := len(line)

        if lenLine == 0 {
            continue
        }
        first := line[0]
        if first == byte('#') || first == byte('!') {
            continue
        }

        sep := []byte{'='}
        index := bytes.Index(line, sep)
        if index < 0 {
            sep = []byte{':'}
        }
        kv := bytes.SplitN(line, sep, 2)
        if kv == nil {
            continue
        }
        lenKV := len(kv)
        if lenKV == 2 {
            p.values[string(bytes.TrimSpace(kv[0]))] = string(bytes.TrimSpace(kv[1]))
        }
    }

    return nil
}

// Write saves the property set to a file. The output will be in "key=value"
// format, with appropriate characters escaped. See Load for more details on
// the file format.
//
// Note: if the property set was loaded from a file, the formatting and
// comments from the original file will not be retained in the output file.
func (p *Properties) Write(w io.Writer) error {
    for k, v := range p.values {
        line := fmt.Sprintf("%s=%s\n", escape(k, true),
            escape(v, false))
        _, err := io.WriteString(w, line)
        if err != nil {
            return err
        }
    }
    return nil
}

// escape returns a string that is safe to use as either a key or value in a
// property file. Whitespace characters, key separators, and comment markers
// should always be escaped.
func escape(s string, key bool) string {

    leading := true
    var buf bytes.Buffer
    for _, ch := range s {
        wasSpace := false
        if ch == '\t' {
            buf.WriteString(`\t`)
        } else if ch == '\n' {
            buf.WriteString(`\n`)
        } else if ch == '\r' {
            buf.WriteString(`\r`)
        } else if ch == '\f' {
            buf.WriteString(`\f`)
        } else if ch == ' ' {
            if key || leading {
                buf.WriteString(`\ `)
                wasSpace = true
            } else {
                buf.WriteRune(ch)
            }
        } else if ch == ':' {
            buf.WriteString(`\:`)
        } else if ch == '=' {
            buf.WriteString(`\=`)
        } else if ch == '#' {
            buf.WriteString(`\#`)
        } else if ch == '!' {
            buf.WriteString(`\!`)
        } else if !unicode.IsPrint(ch) || ch > 126 {
            buf.WriteString(fmt.Sprintf(`\u%04x`, ch))
        } else {
            buf.WriteRune(ch)
        }

        if !wasSpace {
            leading = false
        }
    }
    return buf.String()
}
