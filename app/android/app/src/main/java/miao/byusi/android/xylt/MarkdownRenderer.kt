package miao.byusi.android.xylt

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.LocalContentColor
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.AnnotatedString
import androidx.compose.ui.text.SpanStyle
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.buildAnnotatedString
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextDecoration
import androidx.compose.ui.text.withStyle
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

/**
 * 轻量级 Markdown 渲染器
 *
 * 支持的语法：
 *  - 标题：# / ## / ### / ####
 *  - 粗体：**text**  / __text__
 *  - 斜体：*text*  / _text_
 *  - 删除线：~~text~~
 *  - 行内代码：`code`
 *  - 代码块：```code```
 *  - 引用：> text
 *  - 无序列表：- item / * item
 *  - 有序列表：1. item
 *  - 链接：[text](url)
 *  - 分割线：---
 *  - 段落
 */
object MarkdownRenderer {

    /** Markdown 解析结果：行级节点列表 */
    sealed class Block {
        data class Heading(val level: Int, val text: AnnotatedString) : Block()
        data class Paragraph(val text: AnnotatedString) : Block()
        data class CodeBlock(val code: String) : Block()
        data class Quote(val text: AnnotatedString) : Block()
        data class UnorderedList(val items: List<AnnotatedString>) : Block()
        data class OrderedList(val items: List<AnnotatedString>) : Block()
        data object Divider : Block()
    }

    fun parse(markdown: String): List<Block> {
        val lines = markdown.replace("\r\n", "\n").split("\n")
        val blocks = mutableListOf<Block>()
        val paragraphBuffer = StringBuilder()
        val codeBuffer = StringBuilder()
        val quoteBuffer = StringBuilder()
        val ulBuffer = mutableListOf<AnnotatedString>()
        val olBuffer = mutableListOf<AnnotatedString>()
        var inCodeBlock = false

        fun flushParagraph() {
            if (paragraphBuffer.isNotBlank()) {
                blocks += Block.Paragraph(parseInline(paragraphBuffer.toString().trim()))
                paragraphBuffer.clear()
            }
        }
        fun flushQuote() {
            if (quoteBuffer.isNotBlank()) {
                blocks += Block.Quote(parseInline(quoteBuffer.toString().trim()))
                quoteBuffer.clear()
            }
        }
        fun flushUnorderedList() {
            if (ulBuffer.isNotEmpty()) {
                blocks += Block.UnorderedList(ulBuffer.toList())
                ulBuffer.clear()
            }
        }
        fun flushOrderedList() {
            if (olBuffer.isNotEmpty()) {
                blocks += Block.OrderedList(olBuffer.toList())
                olBuffer.clear()
            }
        }

        for (rawLine in lines) {
            val line = rawLine
            if (line.trimStart().startsWith("```")) {
                if (inCodeBlock) {
                    blocks += Block.CodeBlock(codeBuffer.toString().trimEnd('\n'))
                    codeBuffer.clear()
                    inCodeBlock = false
                } else {
                    flushParagraph(); flushQuote(); flushUnorderedList(); flushOrderedList()
                    inCodeBlock = true
                }
                continue
            }
            if (inCodeBlock) {
                if (codeBuffer.isNotEmpty()) codeBuffer.append('\n')
                codeBuffer.append(line)
                continue
            }

            val trimmed = line.trim()
            if (trimmed.isEmpty()) {
                flushParagraph(); flushQuote(); flushUnorderedList(); flushOrderedList()
                continue
            }

            // 标题
            val headingMatch = Regex("^(#{1,4})\\s+(.+)$").find(trimmed)
            if (headingMatch != null) {
                flushParagraph(); flushQuote(); flushUnorderedList(); flushOrderedList()
                val level = headingMatch.groupValues[1].length
                blocks += Block.Heading(level, parseInline(headingMatch.groupValues[2]))
                continue
            }

            // 分割线
            if (trimmed.matches(Regex("^[-*_]{3,}$"))) {
                flushParagraph(); flushQuote(); flushUnorderedList(); flushOrderedList()
                blocks += Block.Divider
                continue
            }

            // 引用
            if (trimmed.startsWith(">")) {
                flushParagraph(); flushUnorderedList(); flushOrderedList()
                val content = trimmed.removePrefix(">").trimStart()
                if (quoteBuffer.isNotEmpty()) quoteBuffer.append('\n')
                quoteBuffer.append(content)
                continue
            } else {
                flushQuote()
            }

            // 无序列表
            val ulMatch = Regex("^[-*+]\\s+(.+)$").find(trimmed)
            if (ulMatch != null) {
                flushParagraph(); flushOrderedList()
                ulBuffer += parseInline(ulMatch.groupValues[1])
                continue
            } else {
                flushUnorderedList()
            }

            // 有序列表
            val olMatch = Regex("^\\d+\\.\\s+(.+)$").find(trimmed)
            if (olMatch != null) {
                flushParagraph(); flushUnorderedList()
                olBuffer += parseInline(olMatch.groupValues[1])
                continue
            } else {
                flushOrderedList()
            }

            // 段落
            if (paragraphBuffer.isNotEmpty()) paragraphBuffer.append('\n')
            paragraphBuffer.append(trimmed)
        }

        flushParagraph(); flushQuote(); flushUnorderedList(); flushOrderedList()
        if (inCodeBlock && codeBuffer.isNotEmpty()) {
            blocks += Block.CodeBlock(codeBuffer.toString().trimEnd('\n'))
        }
        return blocks
    }

    /** 解析行内元素：粗体、斜体、删除线、行内代码、链接 */
    fun parseInline(text: String): AnnotatedString = buildAnnotatedString {
        var i = 0
        val n = text.length
        while (i < n) {
            // 行内代码 `xxx`
            if (text[i] == '`') {
                val end = text.indexOf('`', i + 1)
                if (end > i) {
                    withStyle(SpanStyle(fontFamily = FontFamily.Monospace, background = Color(0xFFF1F1F1))) {
                        append(text.substring(i + 1, end))
                    }
                    i = end + 1
                    continue
                }
            }
            // 粗体 **xxx** 或 __xxx__
            if (i + 1 < n && text[i] == '*' && text[i + 1] == '*') {
                val end = text.indexOf("**", i + 2)
                if (end > i + 1) {
                    withStyle(SpanStyle(fontWeight = FontWeight.Bold)) {
                        append(text.substring(i + 2, end))
                    }
                    i = end + 2
                    continue
                }
            }
            if (i + 1 < n && text[i] == '_' && text[i + 1] == '_') {
                val end = text.indexOf("__", i + 2)
                if (end > i + 1) {
                    withStyle(SpanStyle(fontWeight = FontWeight.Bold)) {
                        append(text.substring(i + 2, end))
                    }
                    i = end + 2
                    continue
                }
            }
            // 删除线 ~~xxx~~
            if (i + 1 < n && text[i] == '~' && text[i + 1] == '~') {
                val end = text.indexOf("~~", i + 2)
                if (end > i + 1) {
                    withStyle(SpanStyle(textDecoration = TextDecoration.LineThrough)) {
                        append(text.substring(i + 2, end))
                    }
                    i = end + 2
                    continue
                }
            }
            // 斜体 *xxx* 或 _xxx_
            if (text[i] == '*' && i + 1 < n && text[i + 1] != ' ') {
                val end = findClosing(text, i + 1, '*')
                if (end > i) {
                    withStyle(SpanStyle(fontStyle = FontStyle.Italic)) {
                        append(text.substring(i + 1, end))
                    }
                    i = end + 1
                    continue
                }
            }
            if (text[i] == '_' && i + 1 < n && text[i + 1] != ' ') {
                val end = findClosing(text, i + 1, '_')
                if (end > i) {
                    withStyle(SpanStyle(fontStyle = FontStyle.Italic)) {
                        append(text.substring(i + 1, end))
                    }
                    i = end + 1
                    continue
                }
            }
            // 链接 [text](url)
            if (text[i] == '[') {
                val closeBracket = text.indexOf(']', i + 1)
                if (closeBracket > i && closeBracket + 1 < n && text[closeBracket + 1] == '(') {
                    val closeParen = text.indexOf(')', closeBracket + 2)
                    if (closeParen > closeBracket) {
                        val linkText = text.substring(i + 1, closeBracket)
                        val url = text.substring(closeBracket + 2, closeParen)
                        withStyle(
                            SpanStyle(
                                color = Color(0xFF007AFF),
                                textDecoration = TextDecoration.Underline
                            )
                        ) {
                            append(linkText)
                        }
                        // 链接 URL 作为附属信息附加在末尾
                        append(" (")
                        withStyle(SpanStyle(color = Color(0xFF999999), fontSize = 11.sp)) {
                            append(url)
                        }
                        append(")")
                        i = closeParen + 1
                        continue
                    }
                }
            }
            // 转义
            if (text[i] == '\\' && i + 1 < n) {
                append(text[i + 1])
                i += 2
                continue
            }
            append(text[i])
            i++
        }
    }

    private fun findClosing(text: String, from: Int, marker: Char): Int {
        for (i in from until text.length) {
            if (text[i] == marker) return i
        }
        return -1
    }
}

/**
 * 将 Markdown 文本渲染为 Compose UI。
 *
 * 设计目标：纯 Compose，无第三方依赖，支持常见语法。
 */
@Composable
fun MarkdownText(
    markdown: String,
    modifier: Modifier = Modifier,
    color: Color = LocalContentColor.current,
    fontSize: androidx.compose.ui.unit.TextUnit = 15.sp
) {
    val blocks = remember(markdown) { MarkdownRenderer.parse(markdown) }
    Column(modifier = modifier) {
        for (block in blocks) {
            when (block) {
                is MarkdownRenderer.Block.Heading -> {
                    val size = when (block.level) {
                        1 -> (fontSize.value + 10).sp
                        2 -> (fontSize.value + 8).sp
                        3 -> (fontSize.value + 6).sp
                        else -> (fontSize.value + 4).sp
                    }
                    Text(
                        text = block.text,
                        color = color,
                        fontSize = size,
                        fontWeight = FontWeight.Bold,
                        modifier = Modifier.padding(vertical = 4.dp)
                    )
                }
                is MarkdownRenderer.Block.Paragraph -> {
                    Text(
                        text = block.text,
                        color = color,
                        fontSize = fontSize,
                        modifier = Modifier.padding(vertical = 4.dp)
                    )
                }
                is MarkdownRenderer.Block.CodeBlock -> {
                    Box(
                        modifier = Modifier
                            .fillMaxWidth()
                            .clip(RoundedCornerShape(6.dp))
                            .background(Color(0xFFF5F5F5))
                            .padding(10.dp)
                    ) {
                        Text(
                            text = block.code,
                            color = Color(0xFF333333),
                            fontSize = 13.sp,
                            fontFamily = FontFamily.Monospace
                        )
                    }
                    Spacer(modifier = Modifier.height(6.dp))
                }
                is MarkdownRenderer.Block.Quote -> {
                    Box(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(start = 4.dp)
                            .background(Color(0xFFF5F5F5), RoundedCornerShape(4.dp))
                            .padding(10.dp)
                    ) {
                        Text(
                            text = block.text,
                            color = color.copy(alpha = 0.85f),
                            fontSize = fontSize,
                            fontStyle = FontStyle.Italic
                        )
                    }
                    Spacer(modifier = Modifier.height(4.dp))
                }
                is MarkdownRenderer.Block.UnorderedList -> {
                    Column {
                        for (item in block.items) {
                            Row(modifier = Modifier.fillMaxWidth().padding(vertical = 1.dp)) {
                                Text(
                                    text = "•  ",
                                    color = color,
                                    fontSize = fontSize
                                )
                                Text(
                                    text = item,
                                    color = color,
                                    fontSize = fontSize
                                )
                            }
                        }
                    }
                    Spacer(modifier = Modifier.height(4.dp))
                }
                is MarkdownRenderer.Block.OrderedList -> {
                    Column {
                        for ((idx, item) in block.items.withIndex()) {
                            Row(modifier = Modifier.fillMaxWidth().padding(vertical = 1.dp)) {
                                Text(
                                    text = "${idx + 1}.  ",
                                    color = color,
                                    fontSize = fontSize
                                )
                                Text(
                                    text = item,
                                    color = color,
                                    fontSize = fontSize
                                )
                            }
                        }
                    }
                    Spacer(modifier = Modifier.height(4.dp))
                }
                is MarkdownRenderer.Block.Divider -> {
                    Spacer(
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(1.dp)
                            .background(Color(0xFFE0E0E0))
                    )
                    Spacer(modifier = Modifier.height(6.dp))
                }
            }
        }
    }
}
