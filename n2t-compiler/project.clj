(defproject n2t-compiler "0.1.0-SNAPSHOT"
  :description "Compiler for N2T. "
  :url "https://github.com/foggerty/n2t"
  :license {:name "GPL-3.0"
            :url "https://www.gnu.org/licenses/gpl-3.0.en.html"}
  :dependencies [[org.clojure/clojure "1.10.1"]]
  :main ^:skip-aot n2t-compiler.core
  :target-path "target/%s"
  :profiles {:uberjar {:aot :all}})
