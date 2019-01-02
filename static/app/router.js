function config($httpProvider, $compileProvider, $stateProvider, $urlRouterProvider) {
  // ui.router
  $stateProvider
    // base parent state
    .state('parent', {
      views: {
        '': {
          templateUrl: '/static/app/views/layout.html'
        },
        'menu@parent': {
          component: 'menu'
        }
      }
    })
    .state('list', {
      parent: 'parent',
      url: '/',
      data: {
        title: 'List',
      },
      views: {
        'list@parent': {
          component: 'list'
        },
        'content@list': {
          component: 'intro'
        }
      }
    })
    .state('new', {
      parent: 'parent',
      url: '/note/new',
      data: {
        title: 'New note',
      },
      views: {
        'note@parent': {
          component: 'new'
        }
      }
    })
    .state('note', {
      parent: 'parent',
      url: '/note/:uuid',
      data: {},
      resolve: {
        note: function(ApiService, $transition$) {
          return ApiService.getNote($transition$.params().uuid);
        }
      },
      views: {
        'note@parent': {
          component: 'note'
        }
      }
    })
    .state('list_note', {
      parent: 'list',
      url: 'list/:uuid',
      data: {},
      resolve: {
        note: function(ApiService, $transition$) {
          return ApiService.getNote($transition$.params().uuid);
        }

      },
      views: {
        'content@list': {
          component: 'noteRendered'
        }
      }
    })
    ;

  // /emails -> /inbox
  // Automated redirects
  $urlRouterProvider.otherwise('/');

  // batches $digest cycles for $http calls
  // that resolve within 10ms of eachother
  $httpProvider.useApplyAsync(true);

  // BEFORE: <div ng-controller="ListCtrl as list" class="ng-controller ng-binding"></div>
  //         angular.element('.myClass').scope();
  //
  // AFTER:  <div ng-controller="ListCtrl as list"></div>
  //         - Doesn't add unnecessary class names
  //         - Doesn't bind .scope() / .getIsolateScope() data to each element
  $compileProvider.debugInfoEnabled(true);

}

// binding state, param, constants, etc.to root scope
function run($rootScope, $state, $stateParams, APP_SETTINGS) {
  $rootScope.$state = $state;
  $rootScope.$stateParams = $stateParams;
  $rootScope.APP_SETTINGS = APP_SETTINGS;
}

angular
  .module('app')
  .run(run)
  .config(config);
