/* jslint node: true */
"use strict";

function NewCtrl($rootScope, $scope) {
  var vm = this;

  vm.note = {
    id: null,
    content: '',
    title: '',
    created: '',
    updated: ''
  };

  $scope.$watch('vm.note', function(current, original) {
    vm.outputText = marked(current.content);
    $rootScope.$broadcast('currentNote', current);
  }, true);
}

NewCtrl.resolve = {
};

angular.module("app").controller("NewCtrl", NewCtrl);
