package blackhole

import (
    "fmt"
    "github.com/v2fly/v2ray-core/v5/common"
    "github.com/v2fly/v2ray-core/v5/common/buf"
    "github.com/v2fly/v2ray-core/v5/common/serial"
)

const (
    http403response = `HTTP/1.1 403 Forbidden
Connection: close
Cache-Control: max-age=3600, public
Content-Type: text/html
Content-Length: %d

<!DOCTYPE html>
<html>
<head>
    <title>403 Forbidden</title>
</head>
<body>
    <h1>403 Forbidden</h1>
    <p>You do not have permission to access this resource.kevin</p>
</body>
</html>
`
)

// ResponseConfig is the configuration for blackhole responses.
type ResponseConfig interface {
    // WriteTo writes predefined response to the given buffer.
    WriteTo(buf.Writer) int32
}

// WriteTo implements ResponseConfig.WriteTo().
func (*NoneResponse) WriteTo(buf.Writer) int32 { return 0 }

// WriteTo implements ResponseConfig.WriteTo().
func (*HTTPResponse) WriteTo(writer buf.Writer) int32 {
    response := fmt.Sprintf(http403response, len(http403response))
    b := buf.New()
    common.Must2(b.WriteString(response))
    n := b.Len()
    writer.WriteMultiBuffer(buf.MultiBuffer{b})
    return n
}

// GetInternalResponse converts response settings from proto to the internal data structure.
func (c *Config) GetInternalResponse() (ResponseConfig, error) {
    if c.GetResponse() == nil {
        return new(NoneResponse), nil
    }

    config, err := serial.GetInstanceOf(c.GetResponse())
    if err != nil {
        return nil, err
    }
    return config.(ResponseConfig), nil
}
