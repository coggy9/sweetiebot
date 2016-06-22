package sweetiebot

import (
  "github.com/bwmarrin/discordgo"
  "strings"
  "strconv"
  "math/rand"
)

type GiveCommand struct {
}

func (c *GiveCommand) Name() string {
  return "Give";  
}
func (c *GiveCommand) Process(args []string, msg *discordgo.Message) (string, bool) {  
  if len(args) < 1 {
    return "[](/sadbot) `You didn't give me anything!`", false
  }
  if sb.config.MaxBucket == 0 {
    return "```I don't have a bucket right now.```", false 
  }

  arg := ExtraSanitize(strings.Join(args, " "))
  if len(arg) > sb.config.MaxBucketLength {
    return "```That's too big! Give me something smaller!'```", false
  }

  _, ok := sb.config.Collections["bucket"][arg]
  if ok {
    return "```I already have " + arg + "!```", false
  }

  if len(sb.config.Collections["bucket"]) >= sb.config.MaxBucket {
    dropped := BucketDropRandom()
    sb.config.Collections["bucket"][arg] = true
    sb.SaveConfig()
    return "```I dropped " + dropped + " and picked up " + arg + ".```", false
  }

  sb.config.Collections["bucket"][arg] = true
  sb.SaveConfig()
  return "```I picked up " + arg + ".```", false
}
func (c *GiveCommand) Usage() string { 
  return FormatUsage(c, "[arbitrary string]", "Gives sweetie an object. If sweetie is carrying too many things, she will drop one of them at random.") 
}
func (c *GiveCommand) UsageShort() string { return "Gives something to sweetie." }
func (c *GiveCommand) Roles() []string { return []string{} }
func (c *GiveCommand) Channels() []string { return []string{"mylittlebot", "bot-debug"} }

func BucketDropRandom() string {
  index := rand.Intn(len(sb.config.Collections["bucket"]))
  i := 0
  for k, _ := range sb.config.Collections["bucket"] {
    if i == index {
      delete(sb.config.Collections["bucket"], k)
      sb.SaveConfig()
      return k
    }
    i++
  }
  return ""
}

type DropCommand struct {
}

func (c *DropCommand) Name() string {
  return "Drop";  
}

func (c *DropCommand) Process(args []string, msg *discordgo.Message) (string, bool) {  
  if len(sb.config.Collections["bucket"]) == 0 {
    return "```I'm not carrying anything.```", false
  }
  if len(args) < 1 {
    return "```Dropped " + BucketDropRandom() + ".```", false
  }
  arg := strings.Join(args, " ")
  _, ok := sb.config.Collections["bucket"][arg]
  if !ok {
    return "```I don't have " + arg + "!```", false
  }
  delete(sb.config.Collections["bucket"], arg)
  sb.SaveConfig()
  return "```Dropped " + arg + ".```", false
}
func (c *DropCommand) Usage() string { 
  return FormatUsage(c, "[arbitrary string]", "Drops the specified object from sweetie. If no object is given, makes sweetie drop something at random.") 
}
func (c *DropCommand) UsageShort() string { return "Drops something from sweetie's bucket." }
func (c *DropCommand) Roles() []string { return []string{} }
func (c *DropCommand) Channels() []string { return []string{"mylittlebot", "bot-debug"} }


type ListCommand struct {
}

func (c *ListCommand) Name() string {
  return "List";  
}
func (c *ListCommand) Process(args []string, msg *discordgo.Message) (string, bool) {
  things := MapToSlice(sb.config.Collections["bucket"])
  if len(things) == 0 {
    return "```I'm not carrying anything.```", false
  }
  if len(things) == 1 {
    return "```I'm carrying " + things[0] + ".```", false
  }

  return "```I'm carrying " + strings.Join(things[:len(things)-1], ", ") + " and " + things[len(things)-1] + ".```", false
}
func (c *ListCommand) Usage() string { 
  return FormatUsage(c, "", "Lists everything that sweetie has.") 
}
func (c *ListCommand) UsageShort() string { return "Lists everything sweetie has." }
func (c *ListCommand) Roles() []string { return []string{} }
func (c *ListCommand) Channels() []string { return []string{"mylittlebot", "bot-debug"} }

type FightCommand struct {
  monster string
  hp int
}

func (c *FightCommand) Name() string {
  return "Fight";  
}
func (c *FightCommand) Process(args []string, msg *discordgo.Message) (string, bool) {
  things := MapToSlice(sb.config.Collections["bucket"])
  if len(things) == 0 {
    return "```I have nothing to fight with!```", false
  }
  if len(c.monster) > 0 && len(args) > 0 {
    return "I'm already fighting " + c.monster + ", I have to defeat them first!", false
  }
  if len(c.monster) == 0 {
    if len(args) > 0 {
      c.monster = strings.Join(args, " ")
    } else {
      c.monster = sb.db.GetRandomSpeaker()
    }
    c.hp = 10 + rand.Intn(sb.config.MaxFightHP)
    return "```I have engaged " + c.monster + ", who has " + strconv.Itoa(c.hp) + " HP!```", false
  }

  damage := 1 + rand.Intn(sb.config.MaxFightDamage)
  c.hp -= damage
  end := " and deal " + strconv.Itoa(damage) + " damage!"
  monster := c.monster
  if c.hp <= 0 {
    end += " " + monster + " has been defeated!"
    c.monster = ""
  }
  end += "```"
  thing := things[rand.Intn(len(things))]
  switch rand.Intn(7) {
    case 0: return "```I throw " + BucketDropRandom() + " at " + monster + end, false
    case 1: return "```I stab " + monster + " with " + thing + end, false
    case 2: return "```I use " + thing + " on " + monster + end, false
    case 3: return "```I summon " + thing + end, false
    case 4: return "```I cast " + thing + end, false
    case 5: return "```I parry a blow and counterattack with " + thing + end, false
    case 6: return "```I detonate a " + thing + end, false
  }
  return "```Stuff happens" + end, false
}
func (c *FightCommand) Usage() string { 
  return FormatUsage(c, "[name]", "Fights a random pony, or [name] if it is provided.") 
}
func (c *FightCommand) UsageShort() string { return "Fights a random pony." }
func (c *FightCommand) Roles() []string { return []string{} }
func (c *FightCommand) Channels() []string { return []string{"mylittlebot", "bot-debug"} }