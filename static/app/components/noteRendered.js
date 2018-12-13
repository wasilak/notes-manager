/* jslint node: true */
"use strict";

angular.module("app").component("noteRendered", 
  {
    controller: function ($scope, $rootScope, $stateParams, ApiService) {
      var vm = this;

      vm.note = {
        success: true
      };

      ApiService.getNote($stateParams.uuid).then(function(result) {
        vm.note = result;
        $rootScope.$broadcast('currentNote', vm.note);
        vm.inputText = '';
        $rootScope.$broadcast('currentNote', vm.note);

        vm.outputText = marked(vm.note.response.content);
      });

    },
    templateUrl: "/static/app/views/noteRendered.html"
  }
);
