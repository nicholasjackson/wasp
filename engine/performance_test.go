package engine

import (
	"os"
	"runtime"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/wasp/engine/logger"
	"go.uber.org/goleak"

	"runtime/pprof"
)

func perfSetupEngine(t *testing.T, module string, cb *Callbacks) Instance {
	hl := hclog.NewNullLogger()
	hl.SetLevel(hclog.Debug)

	log := logger.New(hl.Info, hl.Debug, hl.Error, hl.Trace)
	e := New(log)

	conf := &PluginConfig{
		Callbacks: cb,
	}

	err := e.RegisterPlugin("test", module, conf)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	inst, err := e.GetInstance("test", "")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	return inst
}

func TestCallsFunctionWithStringParamsMultipleTimes(t *testing.T) {
	t.Skip()
	// check no goroutines are running after the test
	defer goleak.VerifyNone(t)

	//inst := perfSetupEngine(t, "../test_fixtures/go/string_func/module.wasm", nil)
	inst := perfSetupEngine(t, "../_test_fixtures/go/no_imports/module.wasm", nil)

	n := 10000
	for i := 0; i < n; i++ {
		var outString string
		//fmt.Println("############# START ###############")
		//fmt.Println(i)
		//fmt.Println("")

		err := inst.CallFunction("string_func", &outString, longString)
		if err != nil {
			t.Errorf("Failed on iteration %d, error: %s", i, err)
			t.FailNow()
		}

		//fmt.Println("")
		//fmt.Println("############# END #################")
		//fmt.Println("")
		runtime.GC()
	}

	os.Remove("../string_heap.out")
	out, _ := os.Create("../string_heap.out")
	defer out.Close()

	runtime.GC()
	if err := pprof.WriteHeapProfile(out); err != nil {
		t.Fatal("could not write memory profile: ", err)
	}
}

func TestCallsFunctionWithIntParamsMultipleTimes(t *testing.T) {
	t.Skip()
	// check no goroutines are running after the test
	defer goleak.VerifyNone(t)

	inst := perfSetupEngine(t, "../_test_fixtures/go/no_imports/module.wasm", nil)

	n := 10000
	for i := 0; i < n; i++ {
		var outInt int32
		err := inst.CallFunction("int_func", &outInt, 5, 3)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		runtime.GC()
		//fmt.Println("############# START ###############")
		//fmt.Println(i)
		//fmt.Println("")

		//fmt.Println("")
		//fmt.Println("############# END #################")
		//fmt.Println("")
	}

	os.Remove("../int_heap.out")
	out, _ := os.Create("../int_heap.out")
	defer out.Close()

	if err := pprof.WriteHeapProfile(out); err != nil {
		t.Fatal("could not write memory profile: ", err)
	}
}

var longString = `
 They sailed well and the old man soaked his hands in the salt water and tried to keep his head clear. There were high cumulus clouds and enough cirrus above them so that the old man knew the breeze would last all night. The old man looked at the fish constantly to make sure it was true. It was an hour before the first shark hit him.
 The shark was not an accident. He had come up from deep down in the water as the dark cloud of blood had settled and dispersed in the mile deep sea He had come up so fast and absolutely without caution that he broke the surface of the blue water and was in the sun. Then he fell back into the sea and picked up the scent and started swimming on the course the skiff and the fish had taken.
 Sometimes he lost the scent. But he would pick it up again, or have just a trace of it, and he swam fast and hard on the course. He was a very big Mako shark, built to swim as fast as the fastest fish in the sea and everything about him was beautiful except his jaws. His back was as blue as a sword fish's and his belly was silver and his hide was smooth and handsome. He was built as a swordfish except for his huge jaws Which were tight shut now as he swam fast, just under the surface with his high dorsal fin knifing through the water without wavering. Inside the closed double lip of his jaws all of his eight rows of teeth were slanted inwards. They were not the ordinary pyramid-shaped teeth of most sharks. They were shaped like a man's fingers when they are crisped like claws. They were nearly as long as the fingers of the old man and they had razor-sharp cutting edges on both sides. This was a fish built to feed on all the fishes in the sea, that were so fast and strong and well armed that they had no other enemy. Now he speeded up as he smelled the fresher scent and his blue dorsal fin cut the water.
 When the old man saw him coming be knew that this was a shark that had no fear at all and would do exactly what he wished. He prepared the harpoon and made the rope fast while he watched the shark come on. The rope was short as it lacked what he had cut away to lash the fish.
 The old man's head was clear and good now and he was full of resolution, but he had little hope. It was too good to last, he thought. He took one look at the great fish as he watched the shark close in. It might as well have been a dream, he thought. I cannot keep him from hitting me but maybe I can get him. Dentuso , he thought. Bad luck to your mother.
 The shark closed fast astern and when he hit the fish the old man saw his mouth open and his strange eyes and the clicking chop of the teeth as he drove forward in the meat just above the tail. The shark's head was out of the water and his back was coming out and the old man could hear the noise of skin and flesh ripping on the big fish when he rammed the harpoon down onto the shark's head at a spot where the line between his eyes intersected with the line that ran straight back from his nose. There were no such lines. There was only the heavy sharp blue head and the big eyes and the clicking, thrusting all-swallowing jaws. But that was the location of the brain and the old man hit it. He hit it with his blood-mushed hands driving a good harpoon with all his strength. He hit it without hope but with resolution and complete malignancy.
 The shark swung over and the old man saw his eye was not alive and then he swung over once again, wrapping himself in two loops of the rope. The old man knew that he was dead but the shark would not accept it. Then, on his back, with his tail lashing and his jaws clicking, the shark plowed over the water as a speed-boat does. The water was white where his tail beat it and three-quarters of his body was clear above the water when the rope came taut, shivered, and then snapped. The shark lay quietly for a little while on the surface and the old man watched him. Then he went down very slowly.
 "He took about forty pounds," the old man said aloud. He took my harpoon too and all the rope, he thought, and now my fish bleeds again and there will be others.
 He did not like to look at the fish anymore since he had been mutilated. When the fish had been hit it was as though he himself were hit.
 But I killed the shark that hit my fish, he thought. And he was the biggest dentuso that I have ever seen. And God knows that I have seen big ones.
 It was too good to last, he thought. I wish it had been a dream now and that I had never hooked the fish and was alone in bed on the newspapers.
 "But man is not made for defeat," he said. "A man can be destroyed but not defeated." I am sorry that I killed the fish though. Now the bad time is coming and I do not even have the harpoon. The dentuso is cruel and able and strong and intelligent. But I was more intelligent than he was. Perhaps not, he thought Perhaps I was only better armed.
 "Don't think, old man," he said aloud. "Sail on this course and take it when it comes."
 But I must think, he thought. Because it is all I have left. That and baseball. I wonder bow the great DiMaggio would have liked the way I hit him in the brain? It was no great thing, he thought. Any man could do it. But do you think my hands were as great a handicap as the bone spurs? I cannot know. I never had anything wrong with my heel except the time the stingray stung it when I stepped on him when swimming and paralyzed the lower leg and made the unbearable pain.
`
