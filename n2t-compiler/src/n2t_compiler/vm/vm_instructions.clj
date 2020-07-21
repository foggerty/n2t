;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;;;; All of the VM commands put into collections to make lookups easier.

(ns n2t.vm
  (:require [clojure.set :as set]))

(def vm-cmds-arithmetic
  #{:add :sub :neg :eq :gt :lt :and :or :not})

(def vm-cmds-memory
  #{:push :pop})

(def vm-cmds-control
  #{:label :goto :if-goto})

(def vm-cmds-function
  #{:function :call :return})

(def vm-cmds-all
  (set/intersection vm-cmds-arithmetic
                    vm-cmds-control
                    vm-cmds-function
                    vm-cmds-memory))
