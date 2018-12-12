/* jslint node: true */
"use strict";

function NewCtrl($rootScope, $scope) {
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

  $scope.$watch('vm.note', function(current, original) {
    vm.outputText = marked(current.response.content);
    $rootScope.$broadcast('currentNote', current);
  }, true);
}

NewCtrl.resolve = {
};

angular.module("app").controller("NewCtrl", NewCtrl);
