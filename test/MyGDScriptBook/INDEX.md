# Unity.Log("Title")


```c#
[ColorUsage(true, true)] public Color HarmoniousMonsterBubbleColor;
public ColorByColorType BubbleColors;
public Material HarmoniousMonsterMaterial;
public MaterialsByColorType Materials;
public LineRenderer lr_Laser; //控制的lr
public LineRenderer lr_Shadow;
public float duration; //激光消失时间
public float maxLength; //最大距离
```

## [ToRead](GDScript.md) OR [Git](https://github.com)

else if (durationTimer >= duration / 2f && SpellHoverTimer < SpellHoverTime && !shootDone)
{
SpellHoverTimer += Time.deltaTime;
durationTimer -= Time.deltaTime;
hoverRecheckTargetTimer += Time.deltaTime;