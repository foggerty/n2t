;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;;;; Memory management and segment  definitions.

(ns n2t.memory
  (:require [clojure.set :as set]))

(def registers
  "Registers (defined by the assembler) that the VM can access."
  #{:sp :lcl :arg :this :that :R13 :R14 :R15})

(def register-addresses
  {:sp   0
   :lcl  1
   :arg  2
   :this 3
   :that 4
   :temp 5})

(def vm-segments
  #{:argument :local :static :constant :this :that :pointer :temp})

(def segment-address
  "Maps segment to memory address."
  {:stack 256
   :heap  2048})

(def temp
  "Start and end of temp variables."
  {:start 5
   :end   12})

(def general
  "General purpose registers for the VM."
  {:start 13
   :end   15})

(def static
  "Start and end of VM's static variables in RAM."
  {:start 16
   :end   255})

(def stack
  "Start and end of stack in RAM."
  {:start   256
  :end      2047})

(def heap
  "Start and end of heap in RAM."
  {:start   2048
  :end      16383})
