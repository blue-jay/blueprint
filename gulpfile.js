// Modules
var gulp = require('gulp');
var favicon = require ('gulp-real-favicon');
var fs = require('fs');
var runSequence = require('run-sequence');	// Using until gulp v4 is released

// Enviroment variables
var env = JSON.parse(fs.readFileSync('./env.json'))
var folderAsset = env.Asset.Folder;
var folderView = env.View.Folder; 

// Other variables
var faviconData = folderAsset + '/dynamic/favicon/data.json';

// SASS Task
gulp.task('sass', function() {
    var sass = require('gulp-sass');
	var ext = require('gulp-ext-replace');
	gulp.src(folderAsset + '/dynamic/sass/**/*.scss')
        // Available for outputStyle: expanded, nested, compact, compressed
        .pipe(sass({outputStyle: 'expanded'}).on('error', sass.logError))
        .pipe(gulp.dest(folderAsset + '/static/css/'));
    return gulp.src(folderAsset + '/dynamic/sass/**/*.scss')
        // Available for outputStyle: expanded, nested, compact, compressed
        .pipe(sass({outputStyle: 'compressed'}).on('error', sass.logError))
		.pipe(ext('.min.css'))
        .pipe(gulp.dest(folderAsset + '/static/css/'));
});

// JavaScript Task
gulp.task('javascript', function() {
	var concat = require('gulp-concat');
	var minify = require('gulp-minify');
	return gulp.src(folderAsset + '/dynamic/js/*.js')
		.pipe(concat('all.js'))
		.pipe(minify({
			ext:{
			    src:'.js',
			    min:'.min.js'
			}
		}))
		.pipe(gulp.dest(folderAsset + '/static/js/'));
});

// jQuery Task
gulp.task('jquery', function() {
	return gulp.src('node_modules/jquery/dist/jquery.min.*')
		.pipe(gulp.dest(folderAsset + '/static/js/'));
});

// Bootstrap Task
gulp.task('bootstrap', function() {
	gulp.src('node_modules/bootstrap/dist/css/bootstrap-theme.min.*')
		.pipe(gulp.dest(folderAsset + '/static/css/'));
	gulp.src('node_modules/bootstrap/dist/css/bootstrap.min.*')
		.pipe(gulp.dest(folderAsset + '/static/css/'));
	gulp.src('node_modules/bootstrap/dist/fonts/*')
		.pipe(gulp.dest(folderAsset + '/static/fonts/'));
	return gulp.src('node_modules/bootstrap/dist/js/bootstrap.min.js')
		.pipe(gulp.dest(folderAsset + '/static/js/'));
});

// Underscore Task
gulp.task('underscore', function() {
	return gulp.src('node_modules/underscore/underscore-min.*')
		.pipe(gulp.dest(folderAsset + '/static/js/'));
});

// Favicon Generation and Injection Task
gulp.task('favicon', function() {
	runSequence('favicon-generate', 'favicon-inject');
});

// Generate the icons. This task takes a few seconds to complete.
// You should run it at least once to create the icons. Then,
// you should run it whenever RealFaviconGenerator updates its
// package (see the favicon-update task below).
gulp.task('favicon-generate', function(done) {
	var favColor = '#525252';
	favicon.generateFavicon({
		masterPicture: folderAsset + '/dynamic/favicon/logo.png',
		dest: folderAsset + '/static/favicon/',
		iconsPath: '/static/favicon/',
		design: {
			ios: {
				pictureAspect: 'backgroundAndMargin',
				backgroundColor: favColor,
				margin: '14%'
			},
			desktopBrowser: {},
			windows: {
				pictureAspect: 'noChange',
				backgroundColor: favColor,
				onConflict: 'override'
			},
			androidChrome: {
				pictureAspect: 'noChange',
				themeColor: favColor,
				manifest: {
					name: 'Blueprint',
					display: 'browser',
					orientation: 'notSet',
					onConflict: 'override',
					declared: true
				}
			},
			safariPinnedTab: {
				pictureAspect: 'silhouette',
				themeColor: favColor
			}
		},
		settings: {
			scalingAlgorithm: 'Mitchell',
			errorOnImageTooSmall: false
		},
		versioning: {
			paramName: 'v1.0',
			paramValue: '3eepn6WlLO'
		},
		markupFile: faviconData
	}, function() {
		done();
	});
});

// Inject the favicon markups in your HTML pages. You should run
// this task whenever you modify a page. You can keep this task
// as is or refactor your existing HTML pipeline.
gulp.task('favicon-inject', function() {
	return gulp.src([folderView + '/partial/favicon.tmpl'])
		.pipe(favicon.injectFaviconMarkups(JSON.parse(fs.readFileSync(faviconData)).favicon.html_code))
		.pipe(gulp.dest(folderView + '/partial/'));
});

// Check for updates on RealFaviconGenerator (think: Apple has just
// released a new Touch icon along with the latest version of iOS).
// Run this task from time to time. Ideally, make it part of your
// continuous integration system.
gulp.task('favicon-update', function(done) {
	var currentVersion = JSON.parse(fs.readFileSync(faviconData)).version;
	return favicon.checkForUpdates(currentVersion, function(err) {
		if (err) {
			throw err;
		}
	});
});

// Watch
gulp.task('watch', function() {
    gulp.watch(folderAsset + '/dynamic/sass/**/*.scss', ['sass']);
	gulp.watch(folderAsset + '/dynamic/js/*.js', ['javascript']);
});

// Init - every task
gulp.task('init', ['sass', 'javascript', 'jquery', 'bootstrap', 'underscore', 'favicon']);

// Default - only run the tasks that change often
gulp.task('default', ['sass', 'javascript']);