   @24576			// last word of screen memory + 1
   D=A				// set it to D
   @screen-end			// RAM[16]
   M=D				// write 24755 to RAM[16]
(INIT)  
   @SCREEN			// beginning of screen memory
   D=A				// save it to D
   @screen-current		// pointer to current word
   M=D				// assign SCREEN (begining address) to it
   @KBD
   D=M
   @key-pressed
   M=D
   @mask
   M=0
   @MASK-ON
   D;JNE
   @DRAW
   0;JMP
(MASK-ON)   
   @mask			// location of bit 'pattern'
   M=!M				// inverse it so we get all 1's (set to 0 above
(DRAW)   
   @mask			// load bit-mask into D
   D=M
   @screen-current		// write it to the first location in screen
   A=M				// screen-current holds address, so copy it to A
   M=D				// then write mask to whatever location that currently is
   @screen-current		//increment position to next word
   M=M+1			// increment current
   D=M				// save into D for comparison
   @screen-end			// load address of screen-end into A
   D=D-M			// subtract the value in A from D to see if we get a 0
   @WAIT			// if finished, go to wait
   D;JEQ
   @DRAW			// otherwise keep drawing
   0;JMP
(WAIT)
   @KBD
   D=M
   @key-pressed
   D=D-M
   @INIT
   D;JNE
   @WAIT
   0;JMP   
