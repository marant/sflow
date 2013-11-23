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

var mapdata= [];
var features =  [];

function loadData() {
    $.getJSON("map/world-110m.geojson", function(data) {

        //for( var i=0; i<features.length; i++) {
        //    console.log(features[i]);
        //}
        //
        for( var i=0; i<data.features.length; i++) {
            var arrCoords = data.features[i].geometry.coordinates[0];
            var feature =  {coords: []}
            var geometry = new THREE.Geometry();

            for ( var j=0; j<arrCoords.length; j++) {
                var easting = parseFloat(arrCoords[j][0]);
                var northing = parseFloat(arrCoords[j][1]);
                var x = convertToEastingCoordinate(easting)-window.innerWidth/2;
                var y = (convertToNorthingCoordinate(northing)-window.innerHeight/2)*-1;
                feature.coords.push({x: x, y: y});
                geometry.vertices.push(new THREE.Vector3(x, y, 0));

            }
            mapdata.push(geometry);
        }

        init();
    });
}

loadData();

var camera, scene, renderer;
var geometry, material, mesh;
var line;
var camera_z = 1500;
var controls;
var lines = [];

function init() {
    renderer = new THREE.WebGLRenderer();
    renderer.setSize(window.innerWidth, window.innerHeight);
    document.body.appendChild(renderer.domElement);

    camera = new THREE.PerspectiveCamera(45, window.innerWidth / window.innerHeight, 1, 10000); 
    camera.position.set(0, 0, camera_z);
    camera.lookAt(new THREE.Vector3(0, 0, 0));

    window.addEventListener('resize', function() {
        var WIDTH = window.innerWidth;
        var HEIGHT = window.innerHeight;

        renderer.setSize(WIDTH, HEIGHT);
        camera.aspect = WIDTH / HEIGHT;

        camera.updateProjectionMatrix();
    });

    controls = new THREE.OrbitControls(camera, renderer.domElement);
    controls.noRotate = true;

    scene = new THREE.Scene();

    var material = new THREE.LineBasicMaterial({
        color: 0x0000ff
    });

    for( var i=0; i<mapdata.length; i++) {
        var line = new THREE.Line(mapdata[i], material);
        scene.add(line);
    }

    animate();
}

function animate() {
    requestAnimationFrame(animate);

    var material = new THREE.LineBasicMaterial({
        color: 0x0000ff
    });

    controls.update();
    renderer.render(scene, camera);
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

function convertToEastingCoordinate(longitude) {
    var width = window.innerWidth;
    return (longitude+180)*(width/360);
}

function convertToNorthingCoordinate(latitude) {
    var latRad = latitude*(Math.PI/180.0);
    var mercN = Math.log(Math.tan((Math.PI/4)+(latRad/2.0)));
    var width = window.innerWidth;
    var height = window.innerHeight;
    return (height/2)-(width*mercN/(2*Math.PI));
}

