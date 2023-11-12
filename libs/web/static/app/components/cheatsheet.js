/* jslint node: true */
/* global angular */
"use strict";

angular.module("app").component("cheatsheet",
    {
        // eslint-disable-next-line no-unused-vars
        controller: function ($scope, $rootScope, $stateParams, ApiService, $state) {
            var vm = this;
            vm.loader = false;
        },
        templateUrl: "/static/app/views/cheatsheet.html"
    }
);
