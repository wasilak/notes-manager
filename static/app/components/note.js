/* jslint node: true */
"use strict";

angular.module("app").component("note", 
  {
    controller: function ($scope, $rootScope, $stateParams, ApiService) {
      var vm = this;

      vm.note = {
        success: true,
        response: {
          content: ''
        }
      };

      ApiService.getNote($stateParams.uuid).then(function(result) {
        vm.note = result;
        vm.note.edit = true;
        $rootScope.$broadcast('currentNote', vm.note);
      });

      $scope.$watch('$ctrl.note.response.content', function(current, original) {
        vm.errorMessage = false;
        try {
          vm.outputText = marked(current);
        } catch (err) {
          vm.errorMessage = err.message;
        }
      });
    },
    templateUrl: "/static/app/views/note.html"
  }
);
