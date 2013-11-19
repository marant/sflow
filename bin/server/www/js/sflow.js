// WHEE FOR HARD CODING!
var connection = new WebSocket('ws://localhost:12345/data');
connection.binaryType = "arraybuffer";

connection.onopen = function() {
    console.log("connection opened!");
}

connection.onerror = function(error) {
    console.log(error);
}

connection.onmessage = function(blob) {
    var buffer = new Uint8Array(blob.data);
    var packet = parseSFlowPacket(buffer);

    console.log(packet);
}


function parseSFlowPacket(buffer) {
    var sflowPacket = {};

    sflowPacket.srcip = parseIP(buffer, 0);
    sflowPacket.dstip = parseIP(buffer, 4);
    sflowPacket.protocol = buffer[8];
    ports = new(Uint16Array);
    ports[0] = (buffer[9]<<9)+buffer[10]
    ports[1] = (buffer[11]<<9)+buffer[12]

    sflowPacket.srcport = ports[0];
    sflowPacket.dstport = ports[1];

    return sflowPacket;
}

function parseIP(buf, i) {
    var ip = buf[i+3].toString()+"."+buf[i+2].toString()+"."+buf[i+1]+"."+buf[i];
    return ip
}
