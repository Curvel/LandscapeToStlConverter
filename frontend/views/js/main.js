$(document).ready(function() {
    'use strict';

    console.log('main.js ready() fired!');

    $("#form").on("submit", function(){
        event.preventDefault();
        printLandscape.handleGenerateButtonPressed();
    })
});


var printLandscape = (function PrintLandscape() {

    // Private member
    var validateInputFields = function( fields ) {
        if ( !fields || !Array.isArray( fields ) )
            return false;

        var cropping = ['sqr','hex','rnd'];

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

    return {
        // Public member
        handleGenerateButtonPressed: function() {
            var fields = $("#form").serializeArray();
            //JSON.stringify(fields);
            
            if ( !validateInputFields(fields) ) {
                alert('invalid')
                return;
            }

            alert('valid');
        }
    };

})();

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