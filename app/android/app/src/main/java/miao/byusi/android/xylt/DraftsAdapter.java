package miao.byusi.android.xylt;

import android.content.Context;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Button;
import android.widget.TextView;

import androidx.annotation.NonNull;
import androidx.recyclerview.widget.RecyclerView;

import java.util.List;

public class DraftsAdapter extends RecyclerView.Adapter<DraftsAdapter.DraftViewHolder> {

    public interface OnItemClickListener {
        void onItemClick(Article article);
    }

    public interface OnPublishClickListener {
        void onPublishClick(Article article);
    }

    public interface OnDeleteClickListener {
        void onDeleteClick(Article article);
    }

    private Context context;
    private List<Article> drafts;
    private OnItemClickListener itemClickListener;
    private OnPublishClickListener publishClickListener;
    private OnDeleteClickListener deleteClickListener;

    public DraftsAdapter(Context context, List<Article> drafts) {
        this.context = context;
        this.drafts = drafts;
    }

    public void setOnItemClickListener(OnItemClickListener listener) {
        this.itemClickListener = listener;
    }

    public void setOnPublishClickListener(OnPublishClickListener listener) {
        this.publishClickListener = listener;
    }

    public void setOnDeleteClickListener(OnDeleteClickListener listener) {
        this.deleteClickListener = listener;
    }

    @NonNull
    @Override
    public DraftViewHolder onCreateViewHolder(@NonNull ViewGroup parent, int viewType) {
        View view = LayoutInflater.from(context).inflate(R.layout.item_draft, parent, false);
        return new DraftViewHolder(view);
    }

    @Override
    public void onBindViewHolder(@NonNull DraftViewHolder holder, int position) {
        Article draft = drafts.get(position);
        holder.bind(draft);
    }

    @Override
    public int getItemCount() {
        return drafts.size();
    }

    class DraftViewHolder extends RecyclerView.ViewHolder {
        TextView tvTitle;
        TextView tvContent;
        TextView tvDate;
        Button btnPublish;
        Button btnDelete;

        DraftViewHolder(@NonNull View itemView) {
            super(itemView);
            tvTitle = itemView.findViewById(R.id.tv_title);
            tvContent = itemView.findViewById(R.id.tv_content);
            tvDate = itemView.findViewById(R.id.tv_date);
            btnPublish = itemView.findViewById(R.id.btn_publish);
            btnDelete = itemView.findViewById(R.id.btn_delete);
        }

        void bind(Article draft) {
            tvTitle.setText(draft.getTitle().isEmpty() ? "无标题" : draft.getTitle());
            tvContent.setText(draft.getContent());
            tvDate.setText(draft.getCreatedAt());

            itemView.setOnClickListener(v -> {
                if (itemClickListener != null) {
                    itemClickListener.onItemClick(draft);
                }
            });

            btnPublish.setOnClickListener(v -> {
                if (publishClickListener != null) {
                    publishClickListener.onPublishClick(draft);
                }
            });

            btnDelete.setOnClickListener(v -> {
                if (deleteClickListener != null) {
                    deleteClickListener.onDeleteClick(draft);
                }
            });
        }
    }
}
