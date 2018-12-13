/* jslint node: true */
"use strict";

angular.module("app").component("new", 
  {
    controller: function ($scope, $rootScope, $stateParams, ApiService) {
      var vm = this;

      vm.note = {
        success: true,
        response: {
          id: null,
          content: '',
          title: '',
          created: '',
          updated: ''
        }
      };

      $scope.$watch('$ctrl.note', function(current, original) {
        vm.errorMessage = false;
        try {
          vm.outputText = marked(current.response.content);
        } catch (err) {
          vm.errorMessage = err.message;
        }

        $rootScope.$broadcast('currentNote', current);
      }, true);
    },
    templateUrl: "/static/app/views/note.html"
  }
);
