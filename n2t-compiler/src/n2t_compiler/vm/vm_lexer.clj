;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;;;; Parser code for the VM backend.
;;;; Take a directory of .vm files, and output a single collection of VM
;;;; instructions.

(ns lexer
  (:require [clojure.java.io :as io]
            [clojure.string :as str]))

(defn vm-files [dir]
  "Returns all .vm files from dir as Java File objects."
  (filter #(str/ends-with? (str/lower-case (.getName %)) ".vm")
          (.listFiles (io/file dir))))

(defn clean-source [src]
  "Removes comments and empty lines from source file."
  (filter #(not (or (str/starts-with? % "//")
                    (str/blank? %)))
          src))

(defn make-tokens [src]
  "For line in src, returns a map containing :cmd, :arg1 & :arg2."
  (map #(zipmap [:cmd :arg1 :arg2]
                (str/split % #"\s"))
       src))

(defn tokenize [dir]
  "Takes *.vm from dir, and outputs a lazy seq of VM tokens.  Each
  'token' is a map with at least a single key, :cmd.  :arg1
  and :arg2 may also be present."
  (let [files (map #(.getPath %) (vm-files dir))]
    (map
     #(-> (slurp %)
          (str/split-lines)
          (clean-source)
          (make-tokens))
     files)))
