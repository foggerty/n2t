;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;;;; All of the VM commands.

(ns n2t.commands
  (:require [clojure.set :as set]))

(def vm-cmds-arithmetic
  #{:add :sub :neg :eq :gt :lt :and :or :not})

(def vm-cmds-unary
  "List of all VM commands that are unary."
  #{:neg :not})

(def vm-cmds-binary
  "List of all VM commands that are binary operations."
  #{:add :sub :eq :gt :lt :and :not})

(def vm-cmds-memory
  #{:push :pop})

(def vm-cmds-control
  #{:label :goto :if-goto})

(def vm-cmds-function
  #{:function :call :return})

(def vm-cmds-all
  (set/union vm-cmds-arithmetic
             vm-cmds-control
             vm-cmds-function
             vm-cmds-memory))
