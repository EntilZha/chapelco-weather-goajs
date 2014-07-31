'use strict';


// Declare app level module which depends on filters, and services
angular.module('chapelcoWeatherApp', [
  'ngRoute',
	'ui.bootstrap',
  'chapelcoWeatherApp.filters',
  'chapelcoWeatherApp.services',
  'chapelcoWeatherApp.directives',
  'chapelcoWeatherApp.controllers'
]).
config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/', {templateUrl: 'partials/home.html', controller: 'HomeCtrl'});
  $routeProvider.otherwise({redirectTo: '/'});
}]);
