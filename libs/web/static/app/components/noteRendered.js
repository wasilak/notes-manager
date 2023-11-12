/* jslint node: true */
/* jshint -W117 */
/* global angular */
/* global marked */
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

      // eslint-disable-next-line no-unused-vars
      $scope.$watch('$ctrl.note.response', function (current, original) {
        $rootScope.$state.current.data.title = current.title;
        vm.errorMessage = false;
        try {
          vm.outputText = marked.parse(current.content);
        } catch (err) {
          vm.errorMessage = err.message;
        }
      });

      vm.deleteNote = function () {
        var confirmed = confirm("Are you sure?");

        vm.loader = true;

        if (confirmed) {

          // eslint-disable-next-line no-unused-vars
          ApiService.deleteNote(vm.note.response.id).then(function (result) {
            vm.loader = false;
            $rootScope.notifications.push('Note deleted');
            $state.go('list', {}, { reload: true });
          });
        } else {
          vm.loader = false;
        }
      };

    },
    templateUrl: "/static/app/views/noteRendered.html"
  }
);
