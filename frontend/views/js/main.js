// Loaded
$(document).ready(function() {
    'use strict';

    console.log('main.js ready() fired!');

    $('#form').on('submit', function( e ){
        e.preventDefault();
        printLandscape.handleGenerateButtonPressed();
    });

    $('#download').on('click', function( e ){
        printLandscape.downloadLastFile();
    });
});

// App
var printLandscape = (function PrintLandscape() {

    const status = {
        IDLE: 'idle',
        IN_PROGRESS: 'in_progress',
        SUCCESS: 'success',
        FAILED: 'failed'
    }

    var socket;
    var lastFileUrl;
    var stlViewer;

    // Private member
    var validateInputFields = function( fields ) {
        if ( !fields || !Array.isArray( fields ) )
            return false;

        var cropping = ['sqr','hex','rnd'];

        // TODO: Add heightFactor and length, ...
        for ( field of fields ) {
            if ( !field.name || !field.value )
                return false;

            switch (field.name) {
                case 'cropping':
                        if (cropping.includes(field.value))
                            continue;
                    break;

                case 'northEastLat':
                case 'southWestLat':
                    if ( !isNaN(field.value) || field.value >= -90 || field.value <= 90 )
                        continue;
                    break;

                case 'northEastLng':
                case 'southWestLng':
                    if ( !isNaN(field.value) || field.value >= -180 || field.value <= 180 )
                        continue;
                    break;
            
                default:
                    break;
            }

            return false;
        }

        return true;
    };

    var setStatus = function( current ) {
        switch( current ) {
            case status.IDLE:
                $('#submit').prop('disabled', false);
                $('#progressbar').attr('hidden','');
                $('#alert').attr('hidden','');
                resetPercentage();
                break;
            case status.IN_PROGRESS:
                $('#submit').prop('disabled', true);
                $('#progressbar').removeAttr('hidden');
                scrollToBottom();
                break;
            case status.SUCCESS:
                $('#submit').prop('disabled', false);
                $('#download').removeAttr('hidden');
                break;
            case status.FAILED:
                $('#submit').prop('disabled', false);
                $('#progressbar').attr('hidden','');
                $('#alert').removeAttr('hidden');
                resetPercentage();
                break;
        }
    };

    var submitRequest = function( data ) {       
        socket = io.connect('http://127.0.0.1:4321', {
            reconnection: false
        });
        socket.on('error', console.error.bind(console));
        socket.on('message', console.log.bind(console));

        socket.on('convertUpdate', function (data) {
            console.log(data);

            var percentage = data.split(';');
            convertUpdate(percentage[0], percentage[1], percentage[2]);
        });
        socket.on('convertSuccess', function (data) {
            console.log(data);
            setStatus(status.SUCCESS);

            lastFileUrl = data;
            loadPreview(lastFileUrl);

            socket.disconnect();
        });
        socket.on('convertFailed', function (data) {
            console.log(data);
            setStatus(status.FAILED);

            socket.disconnect();
        });

        socket.emit('requestConvert', {fields: data, id: 'unused'});
    };
    var loadPreview = function( pathToFile ) {
        $('.preview-section').removeAttr('hidden');
        $('#stlViewer').find('canvas:first').remove();

        // Should be: /files/file.stl
        stlViewer = new StlViewer(document.getElementById("stlViewer"), { models: [ {id: 0, filename: pathToFile} ] });
        scrollToBottom();
    };
    var scrollToBottom = function() {
        $('html, body').animate({scrollTop:$(document).height()}, 'slow');
    };

    // Progress
    var convertUpdate = function( percentage1, percentage2, percentage3 ) {

        setPartPercentage('#progressbarStep1', 'Generate height map...', percentage1);
        setPartPercentage('#progressbarStep2', 'Triangulation...', percentage2);
        setPartPercentage('#progressbarStep3', 'STL generation...', percentage3);
    };
    var resetPercentage = function() {
        setPartPercentage('#progressbarStep1');
        setPartPercentage('#progressbarStep2');
        setPartPercentage('#progressbarStep3');
    };
    var setPartPercentage = function( elementSelector, title, totalPercentage ) {
        if ( totalPercentage && totalPercentage > 0 ) {

            var partPerventage = Math.round(totalPercentage / 3);

            // Workaround for third bar
            if (partPerventage == 33 && elementSelector == '#progressbarStep3')
                partPerventage = 34;

            console.log(partPerventage + ' of ' + elementSelector);
            $(elementSelector).text(title + ' ' + totalPercentage + '%');
            $(elementSelector).width(partPerventage + "%");
        } else {
            $(elementSelector).text('');
            $(elementSelector).width('0%');
        }
    };

    var tryToDownloadLastFile = function() {
        if ( lastFileUrl ) {
            window.open(document.location.origin + lastFileUrl, '_blank');
            //alert('Redirect to: ' + document.location.origin + lastFileUrl);
        } else {
            alert('None file generated yet.');
        }
    };

    return {
        // Public member
        handleGenerateButtonPressed: function() {
            setStatus(status.IDLE);

            var fields = $('#form').serializeArray();
            //JSON.stringify(fields);
            
            // TODO: Implement
            /* if ( !validateInputFields(fields) ) {
                alert('Validation failed.')
                return;
            } */

            console.log('Validation passed.');

            setStatus(status.IN_PROGRESS);
            submitRequest(fields);
        },
        downloadLastFile: function() {
            tryToDownloadLastFile();
        }
    };

})();

// Google Maps
var mapsIntegration = (function MapsIntegraion() {
    var map;
    return {
        initMap: function() {
            map = new google.maps.Map(document.getElementById('map'), {
                center: {lat: 49.0113, lng: 8.4192},
                zoom: 8
            });
        
            google.maps.event.addListener(map, 'bounds_changed', function() {
                var bounds =  map.getBounds();
                var ne = bounds.getNorthEast(); // >^
                var sw = bounds.getSouthWest(); // <v
        
                // console.log(ne + ' | ' + sw);
                $('input[name="northEastLat"]').val(ne.lat());
                $('input[name="northEastLng"]').val(ne.lng());
        
                $('input[name="southWestLat"]').val(sw.lat());
                $('input[name="southWestLng"]').val(sw.lng());
            });
        }
    }
})();