;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;;;; Parser code for the VM backend.
;;;; Take a directory of .vm files, and output a single collection of VM
;;;; instructions.

(ns n2t.lexer
  (:require [n2t.vm :refer :all]
            [clojure.java.io :as io]
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

(defn map-cmd [token]
  "Returns token with :cmd string mapped to a valid VM command."
  (let [cmd (keyword (:cmd token))]
    (if (contains? vm-cmds-all cmd)
      (assoc token :cmd cmd)
      (throw (Exception. "Invalid token name, learn how string interpolation works in Clojure.")))))

(defn map-arg1 [token]
  "Repalce :arg1 with a keyword, if a matching one can be found (could
  be a user-defined function name, so cannot just apply the keyword
  function)."
  (let [arg (keyword (:arg1 token))]
    (if (contains? vm-segments arg)
      (assoc token :arg1 arg)
      token)))

(defn map-values [tokens]
  "Expects a seq of tokens, and will map the values of :cmd and :arg1
  to appropriate keywords."
  (map #(-> (map-cmd %)
            (map-arg1))
       tokens))

(defn tokenize [dir]
  "Takes *.vm from dir, and outputs a lazy seq of VM tokens.  Each
  'token' is a map with at least a single key, :cmd.  :arg1 and :arg2
  may also be present."
  (let [files (map #(.getPath %) (vm-files dir))]
    (flatten
     (map
      #(-> (slurp %)
           (str/split-lines)
           (clean-source)
           (make-tokens)
           (map-values))
      files))))
