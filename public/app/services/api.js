function ApiService($http, API, $rootScope, APP_SETTINGS) {
  var getNote = function(uuid) {

    var url = API.urls.note.replace("{{uuid}}", uuid);

    return $http({
        cache: false,
        url: url,
        method: 'GET',
        params: {}
      })
      .then(function(response) {
        return response.data;
      },
      function(response) {
        console.error('Getting note failed!');
        console.error(response);
        return {
          success: false,
          error: response
        };
      }
    );
  };

  var getList = function() {

    var url = API.urls.list;

    return $http({
        cache: false,
        url: url,
        method: 'GET',
        params: {}
      })
      .then(function(response) {
        return response.data;
      },
      function(response) {
        console.error('Getting notes list failed!');
        console.error(response);
        return {
          success: false,
          error: response
        };
      }
    );
  };

  var saveNote = function(note) {

    var url = API.urls.note.replace("{{uuid}}", note.id);

    return $http({
        cache: false,
        url: url,
        method: 'POST',
        data: {
          note: note
        }
      })
      .then(function(response) {
        return response.data;
      },
      function(response) {
        console.error('Saving note failed!');
        console.error(response);
        return {
          success: false,
          error: response
        };
      }
    );
  };

  return {
    getNote: getNote,
    saveNote: saveNote,
    getList: getList
  };
}

angular.module("app").factory("ApiService", ApiService);
