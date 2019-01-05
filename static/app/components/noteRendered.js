/* jslint node: true */
/* jshint -W117 */
"use strict";

angular.module("app").component("noteRendered", 
  {
    bindings: {
      note: '<'
    },
    controller: function ($scope, $rootScope, $stateParams, ApiService, $state) {
      var vm = this;

      vm.$state = $state;

      vm.loader = false;

      $scope.$watch('$ctrl.note.response', function(current, original) {
        $rootScope.$state.current.data.title = current.title;
        vm.errorMessage = false;
        try {
          vm.outputText = marked(current.content);
        } catch (err) {
          vm.errorMessage = err.message;
        }
      });

      vm.deleteNote = function() {
        var confirmed = confirm("Are you sure?");

        vm.loader = true;

        if (confirmed) {
          ApiService.deleteNote(vm.note.response.id).then(function(result) {
              vm.loader = false;
              $rootScope.notifications.push('Note deleted');
              $state.go('list', {}, {reload: true});
          });
        } else {
          vm.loader = false;
        }
      };

    },
    templateUrl: "/static/app/views/noteRendered.html"
  }
);
