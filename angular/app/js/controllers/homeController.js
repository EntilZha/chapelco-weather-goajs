angular.module('chapelcoWeatherApp.controllers').controller('HomeCtrl', ['$scope', '$http', function($scope, $http) {
	$scope.dataLoaded = false;
	$http.get('api/weather/current').success(function(data) {
		$scope.currentWeather = data;
		$scope.dataLoaded = true;
	});
}]);
