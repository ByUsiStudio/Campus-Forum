package miao.byusi.android.xylt;

import android.content.Context;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ImageView;
import android.widget.TextView;

import androidx.annotation.NonNull;
import androidx.recyclerview.widget.RecyclerView;

import com.bumptech.glide.Glide;

import java.util.List;

public class ArticleAdapter extends RecyclerView.Adapter<ArticleAdapter.ArticleViewHolder> {

    private Context context;
    private List<Article> articleList;
    private OnItemClickListener listener;

    public interface OnItemClickListener {
        void onItemClick(Article article);
    }

    public ArticleAdapter(Context context, List<Article> articleList) {
        this.context = context;
        this.articleList = articleList;
    }

    public void setOnItemClickListener(OnItemClickListener listener) {
        this.listener = listener;
    }

    @NonNull
    @Override
    public ArticleViewHolder onCreateViewHolder(@NonNull ViewGroup parent, int viewType) {
        View view = LayoutInflater.from(context).inflate(R.layout.item_article, parent, false);
        return new ArticleViewHolder(view);
    }

    @Override
    public void onBindViewHolder(@NonNull ArticleViewHolder holder, int position) {
        Article article = articleList.get(position);
        holder.tvTitle.setText(article.getTitle());
        holder.tvContent.setText(article.getContent());
        holder.tvAuthor.setText(article.getAuthor());
        holder.tvViews.setText(String.valueOf(article.getViewCount()));
        holder.tvLikes.setText(String.valueOf(article.getLikeCount()));
        holder.tvComments.setText(String.valueOf(article.getCommentCount()));
        holder.tvDate.setText(article.getCreatedAt());

        Glide.with(context).load(article.getAvatarUrl()).into(holder.ivAvatar);

        if (article.isTop()) {
            holder.tvTopBadge.setVisibility(View.VISIBLE);
        } else {
            holder.tvTopBadge.setVisibility(View.GONE);
        }

        holder.itemView.setOnClickListener(v -> {
            if (listener != null) {
                listener.onItemClick(article);
            }
        });
    }

    @Override
    public int getItemCount() {
        return articleList.size();
    }

    public static class ArticleViewHolder extends RecyclerView.ViewHolder {
        ImageView ivAvatar;
        TextView tvTitle;
        TextView tvContent;
        TextView tvAuthor;
        TextView tvViews;
        TextView tvLikes;
        TextView tvComments;
        TextView tvDate;
        TextView tvTopBadge;

        public ArticleViewHolder(@NonNull View itemView) {
            super(itemView);
            ivAvatar = itemView.findViewById(R.id.iv_avatar);
            tvTitle = itemView.findViewById(R.id.tv_title);
            tvContent = itemView.findViewById(R.id.tv_content);
            tvAuthor = itemView.findViewById(R.id.tv_author);
            tvViews = itemView.findViewById(R.id.tv_views);
            tvLikes = itemView.findViewById(R.id.tv_likes);
            tvComments = itemView.findViewById(R.id.tv_comments);
            tvDate = itemView.findViewById(R.id.tv_date);
            tvTopBadge = itemView.findViewById(R.id.tv_top_badge);
        }
    }
}
