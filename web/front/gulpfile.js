/* global require */
'use strict';

const gulp = require('gulp');
const run = require('run-sequence');
const util = require('gulp-util');
const del = require('del');

// styles
const sass = require('gulp-sass');
const plumber = require('gulp-plumber');
const postcss = require('gulp-postcss');
const csso = require('gulp-csso');
const autoprefixer = require('autoprefixer');
// images
const svgmin = require('gulp-svgmin');
// icons
// vue & js
const browserify = require('browserify');
const aliasify = require('aliasify');
const buffer = require('vinyl-buffer');
const source = require('vinyl-source-stream');
const uglify = require('gulp-uglify');
const babelify = require('babelify');
const sourcemaps = require('gulp-sourcemaps');
// docker
const envify = require('envify/custom');

const vars = {
  backendPath: '../static/',
  production: !!util.env.production
};

const config = {
  style: {
    input: 'scss/**/*.{scss,sass}',
    output: vars.backendPath + '/css/',
    params: {
      outputStyle: 'compress',
    }
  },
  images: {
    input: ['!img/icons', 'img/**/*.{png,jpg,svg,gif}'],
    output: vars.backendPath + '/img',
  },
  icons: {
    input: 'img/icons/*.svg',
    output: vars.backendPath + 'templates/',
  },
  fonts: {
    input: 'fonts/**/*.{woff,woff2,eot,ttf}',
    output: vars.backendPath + '/fonts',
  },
  scripts: {
    watchPath: ['js/**/*.js', 'js/*.{js}', 'js/**/*.vue'],
    vuePath: ['js/main.js'],
    othersPaths: ['js/*.js'],
    output: vars.backendPath + '/js'
  },
  other: {
    input: 'other/**',
    output: vars.backendPath + '/'
  }
}

gulp.task('style', function() {
  return gulp.src(config.style.input)
    .pipe(plumber())
    .pipe(sass(config.style.params))
    .pipe(postcss([
      autoprefixer({browsers: [
        'ie > 9',
        '> 1%',
        'last 5 versions'
      ]})
    ]))
    .pipe(gulp.dest(config.style.output))
    .pipe(csso());
});

gulp.task('images', function() {
  return gulp.src(config.images.input)
    .pipe(gulp.dest(config.images.output));
});

gulp.task('icons', function() {
  return gulp
    .src(config.icons.input)
    .pipe(svgmin(function (file) {
      return {
        plugins: [{
          cleanupIDs: {
            minify: false
          }
        }]
      };
    }));
});

gulp.task('fonts', function() {
  return gulp.src(config.fonts.input).pipe(gulp.dest(config.fonts.output));
});

gulp.task('browserify', function() {
  const b = browserify({
    entries: config.scripts.vuePath,
    transform: [babelify],
  });

  b.transform(aliasify, {
    global: true,
    aliases: {
      'react': 'preact-compat',
      'react-dom': 'preact-compat'
    }
  });

  b.transform(
    // Порядок необходим для обработки файлов node_modules
    { global: true },
    envify({ NODE_ENV: process.env.NODE_ENV })
  );

  return b.bundle()
    .pipe(source('app.js'))
    .pipe(buffer())
  // If --production flag is set, uglify, do not create sourcemaps
    .pipe(vars.production? uglify() : util.noop())
  // for debugging
    .pipe(vars.production? util.noop() : sourcemaps.init({ loadMaps: true }))
    .pipe(vars.production? util.noop() : sourcemaps.write('./'))
    .pipe(gulp.dest(config.scripts.output));
});

gulp.task('other', function() {
  return gulp.src(config.other.input)
    .pipe(gulp.dest(config.other.output));
});

gulp.task('clean', function() {
  return del(vars.backendPath, {force: true});
});

gulp.task('dist', gulp.series(
  'style', 'icons', 'images', 'browserify', 'fonts', 'other'
));

gulp.task('watch', function() {
  gulp.watch(config.style.input, ['style']);
  gulp.watch(config.images.input, ['images']);
  gulp.watch(config.icons.input, ['icons']);
  gulp.watch(config.fonts.input, ['fonts']);
  gulp.watch(config.scripts.input, ['scripts']);
});
