var currentPage = 1;
const pageSize = 5;
var totalPage

const url = "127.0.0.1";

var commit = document.getElementById("commit");
var username = document.getElementById("userName");
var commentcontent = document.getElementById("commentContent");


window.onload = function ()
{
    render();
    commit = document.getElementById("commit");
    username = document.getElementById("userName");
    commentcontent = document.getElementById("commentContent");
}

function render()
{
    var myRequest = new Request(`/comment/get?page=${currentPage}&size=${pageSize}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    })

    fetch(myRequest)
        .then(response => response.json())
        .then(res =>
        {
            console.log(res)
            show_comments(res.data.comments);

            console.log("rendering page number!")
            totalPage = Math.ceil(res.data.total / pageSize);
            if (totalPage == 0) totalPage = 1;
            console.log(res.data.total, pageSize, totalPage);
            var page = document.getElementById("pagex");
            while (page.firstChild)
            {
                page.removeChild(page.firstChild);
            }
            var pagetext = document.createElement("span");
            pagetext.innerHTML = `<strong>${currentPage}/${totalPage}</strong>`;
            page.appendChild(pagetext);
        })
}

// 添加评论
commit.addEventListener("click", () =>
{
    if (username.value == "" || commentcontent.value == "")
    {
        alert("请输入用户名和评论内容😡");
        return;
    }
    var myRequest = new Request(`/comment/add`, {
        method: "POST",
        body: JSON.stringify({
            "name": username.value,
            "content": commentcontent.value
        }),
        headers: {
            "Content-Type": "application/json"
        }
    })
    fetch(myRequest)
        .then(response => response.json())
        .then(data =>
        {
            console.log(data);
        })
        .then(render)
    clear_input();
})



function show_comments(comments)
{
    var div = document.getElementById("comments");
    while (div.firstChild)
    {
        div.removeChild(div.firstChild);
    }

    for (var i = 0; i < comments.length; i++)
    {
        var comment = comments[i];
        var commentElement = document.createElement("div");
        commentElement.className = "commentsZone";
        commentElement.innerHTML = `
        <div class="userName_commentContent">
            <h3>${comment.name}</h3>
        </div>

        <div class="userName_commentContent">
            <h4>${comment.content}</h4>
        </div>

        <div class="buttonLine">
            <button class="buttonStyle" id="delete${comment.id}">删除</button>
        </div>
        `;
        div.appendChild(commentElement);
    }

    var buttons = document.querySelectorAll('[id^="delete"]');
    buttons.forEach(button =>
    {
        button.addEventListener("click", () =>
        {
            var id = button.id.replace("delete", "");
            delete_comment(id);
        })
    })
}

function delete_comment(id)
{
    var myRequest = new Request(`/comment/delete?id=${id}`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        }
    })
    fetch(myRequest).then(() => render())
}

function clear_input()
{
    username.value = "";
    commentContent.value = "";
}


var PREV = document.getElementById("previous");
var NEXT = document.getElementById("next");

PREV.addEventListener("click", () =>
{
    if (currentPage >= 2)
    {
        currentPage--;
        render();
    }
})

NEXT.addEventListener("click", () =>
{
    if (currentPage < totalPage)
    {
        currentPage++;
        render();
    }
})